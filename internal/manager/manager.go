package manager

import (
	"context"
	"fmt"
	"match_evaluator/evaluator"
	"match_evaluator/internal/db"
	match_evaluator "match_evaluator/proto"
	"strconv"
	"sync"
	"time"

	match_process "github.com/askldfhjg/match_apis/match_process/proto"

	"github.com/micro/micro/v3/service/broker"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/metrics"
	"google.golang.org/protobuf/proto"
)

func NewManager(opts ...evaluator.EvaluatorOption) evaluator.Manager {
	m := &defaultMgr{
		exited1:        make(chan struct{}, 1),
		exited2:        make(chan struct{}, 1),
		exited3:        make(chan struct{}, 1),
		msgList:        make(chan interface{}, 1000),
		msgBackList:    make(chan string, 100),
		resultChannel1: make(chan []*match_evaluator.MatchDetail, 1000),
		resultChannel2: make(chan *match_evaluator.MatchDetail, 100000),
		channelMap:     make(map[string]chan interface{}, 100),
		versionMap:     &sync.Map{},
	}
	for _, o := range opts {
		o(&m.opts)
	}
	return m
}

type defaultMgr struct {
	opts           evaluator.EvaluatorOptions
	exited1        chan struct{}
	exited2        chan struct{}
	exited3        chan struct{}
	msgList        chan interface{}
	msgBackList    chan string
	resultChannel1 chan []*match_evaluator.MatchDetail
	resultChannel2 chan *match_evaluator.MatchDetail
	channelMap     map[string]chan interface{}
	versionMap     *sync.Map
}

func (m *defaultMgr) Start() error {
	err := m.SubPoolVersion()
	if err != nil {
		return err
	}
	go m.resultPublishLoop1()
	go m.resultPublishLoop2()
	go m.loop()
	return nil
}

func (m *defaultMgr) Stop() error {
	close(m.exited1)
	close(m.exited2)
	close(m.exited3)
	return nil
}

func (m *defaultMgr) SubPoolVersion() error {
	_, err := broker.Subscribe("pool_version", m.handlerMsg)
	return err
}

func (m *defaultMgr) handlerMsg(raw *broker.Message) error {
	msg := &match_process.PoolVersionMsg{}
	err := proto.Unmarshal(raw.Body, msg)
	if err != nil {
		return err
	}
	logger.Infof("handlerMsg %+v", msg)
	m.versionMap.Store(fmt.Sprintf("%s:%d", msg.GameId, msg.SubType), msg.Version)
	return nil
}

func (m *defaultMgr) loop() {
	for {
		select {
		case <-m.exited1:
			return
		case req := <-m.msgList:
			var gameId string
			var subType int64
			var nowVersion int64
			var evalGroupId string
			switch v := req.(type) {
			case *match_evaluator.ToEvalReadyReq:
				gameId = v.GameId
				subType = v.SubType
				nowVersion = v.Version
				evalGroupId = v.EvalGroupId
			case *match_evaluator.ToEvalReq:
				gameId = v.GameId
				subType = v.SubType
				nowVersion = v.Version
				evalGroupId = v.EvalGroupId
			}
			if len(gameId) > 0 {
				key := fmt.Sprintf("%s:%d", gameId, subType)
				version, err := m.getPoolVersion(key)
				if err != nil {
					logger.Errorf("get version1 from redis have err %s", err.Error())
				} else if version == nowVersion {
					channel, ok := m.channelMap[evalGroupId]
					if !ok {
						channel = make(chan interface{}, 10)
						m.channelMap[evalGroupId] = channel
						go m.processEval(channel, evalGroupId)
					}
					channel <- req
				}
			}
		case groupId := <-m.msgBackList:
			delete(m.channelMap, groupId)
		}
	}
}

func (m *defaultMgr) resultPublishLoop1() {
	for {
		select {
		case <-m.exited2:
			return
		case reqs := <-m.resultChannel1:
			for _, bb := range reqs {
				m.resultChannel2 <- bb
			}
		}
	}
}

func (m *defaultMgr) resultPublishLoop2() {
	for {
		select {
		case <-m.exited3:
			return
		case req := <-m.resultChannel2:
			err := m.publishResult(req)
			if err != nil {
				//logger.Errorf("PublishResult have err %s", err.Error())
			} else {
				//logger.Info("fffffff")
			}
		}
	}
}

func (m *defaultMgr) AddMsg(req interface{}) {
	m.msgList <- req
}

type detailResult struct {
	List               []*match_evaluator.MatchDetail
	MissIndex          map[int]bool
	EvalGroupSubId     int64
	EvalGroupTaskCount int64
	Version            int64
	Key                string
	StartTime          int64
	RunTime            int64
	GameId             string
	SubType            int64
	EvalStartTime      int64
}

func (m *defaultMgr) processEval(chnl chan interface{}, gId string) {
	getIndex1 := map[int]bool{}
	getIndex2 := map[int]bool{}
	getIndex3 := map[int]bool{}
	redayOk := false
	groupId := gId
	channel := chnl

	inTeam := make(map[string]bool)
	timer := time.NewTimer(3 * time.Second)
	msgChannel := make(chan *detailResult, 100)
	tmpList := make([]*detailResult, 0, 32)
	var startTime int64
	go func() {
	FOR:
		for {
			select {
			case <-timer.C:
				timer.Stop()
				m.msgBackList <- groupId
				break FOR
			case req := <-msgChannel:
				if redayOk {
					flag := getIndex3[int(req.EvalGroupSubId)]
					if !flag {
						getIndex3[int(req.EvalGroupSubId)] = true
						m.Rem2(req)
						m.metricsEach(req)
					}
					if len(getIndex3) == int(req.EvalGroupTaskCount) {
						timer.Stop()
						m.msgBackList <- groupId
						logger.Infof("delete result timer %v", time.Now().UnixNano()/1e6)
						m.metrics(req)
					}
				} else {
					tmpList = append(tmpList, req)
				}
			case req := <-channel:
				switch v := req.(type) {
				case *match_evaluator.ToEvalReadyReq:
					flag := getIndex1[int(v.EvalGroupSubId)]
					if !flag {
						getIndex1[int(v.EvalGroupSubId)] = true
					}
					if len(getIndex1) == int(v.EvalGroupTaskCount) && !redayOk {
						logger.Infof("ready ok %v", time.Now().UnixNano()/1e6)
						redayOk = true
						if len(tmpList) > 0 {
							for _, v := range tmpList {
								msgChannel <- v
							}
							tmpList = tmpList[:0]
						}
					}
				case *match_evaluator.ToEvalReq:
					flag := getIndex2[int(v.EvalGroupSubId)]
					if !flag {
						key := fmt.Sprintf("%s:%d", v.GameId, v.SubType)
						getIndex2[int(v.EvalGroupSubId)] = true
						if startTime < 0 {
							startTime = time.Now().UnixNano() / 1e6
						}
						m.Eval2(v, key, startTime, msgChannel, &inTeam)
					}
					if len(getIndex2) == int(v.EvalGroupTaskCount) {
						logger.Infof("result timer %v", time.Now().UnixNano()/1e6)
					}
				}
			}
		}
		logger.Infof("%s end", gId)
	}()
}

func (m *defaultMgr) metrics(req *detailResult) {
	now := time.Now().UnixNano() / 1e6
	tags := metrics.Tags{
		"gameId":  req.GameId,
		"subType": strconv.FormatInt(req.SubType, 10),
	}
	err := metrics.Timing("match.fulltime", time.Duration(req.StartTime-now)*time.Millisecond, tags)
	if err != nil {
		logger.Errorf("match.fulltime:err:%s", err.Error())
	}

	err = metrics.Timing("match.evaltime", time.Duration(req.EvalStartTime-now)*time.Millisecond, tags)
	if err != nil {
		logger.Errorf("match.evaltime:err:%s", err.Error())
	}
}

func (m *defaultMgr) metricsEach(req *detailResult) {
	tags := metrics.Tags{
		"gameId":  req.GameId,
		"subType": strconv.FormatInt(req.SubType, 10),
	}

	err := metrics.Timing("match.evaltime", time.Duration(req.RunTime)*time.Millisecond, tags)
	if err != nil {
		logger.Errorf("match.evaltime:err:%s", err.Error())
	}
}

// func (m *defaultMgr) Eval(list []*match_evaluator.MatchDetail, groupId string, version int64, keyy string, inTeam *map[string]bool) {
// 	// if !pass {
// 	// 	sort.Slice(list, func(i, j int) bool {
// 	// 		return list[i].Score < list[j].Score
// 	// 	})
// 	// }
// 	// inTeam := map[string]int{}
// 	ctx := context.Background()
// 	count := 0
// 	deleteList := make([]string, 0, deleteGroupCount)
// 	retDetail := make([]*match_evaluator.MatchDetail, 0, deleteGroupCount)
// 	for _, detail := range list {
// 		have := false
// 		for _, ply := range detail.Ids {
// 			_, ok := (*inTeam)[ply]
// 			have = have || ok
// 		}

// 		if !have {
// 			for _, ply := range detail.Ids {
// 				(*inTeam)[ply] = true
// 			}
// 			nowV, err := m.getPoolVersion(keyy)
// 			if err == nil && nowV == version {
// 				deleteList = append(deleteList, detail.Ids...)
// 				retDetail = append(retDetail, detail)
// 				if len(deleteList) >= deleteGroupCount {
// 					delCount, err := db.Default.RemoveTokens(ctx, deleteList, keyy)
// 					if err == nil {
// 						count += len(retDetail)
// 						m.resultChannel1 <- retDetail
// 						if delCount != len(deleteList) {
// 							logger.Errorf("eval delCount have err %d %d", delCount, len(deleteList))
// 						}
// 					} else {
// 						logger.Errorf("RemoveTokens have err %s", err.Error())
// 						break
// 					}
// 					deleteList = make([]string, 0, deleteGroupCount)
// 					retDetail = make([]*match_evaluator.MatchDetail, 0, deleteGroupCount)
// 				}
// 			} else {
// 				logger.Infof("result Count version break")
// 				break
// 			}
// 		}
// 	}
// 	if len(deleteList) > 0 {
// 		delCount, err := db.Default.RemoveTokens(ctx, deleteList, keyy)
// 		if err == nil {
// 			count += len(retDetail)
// 			m.resultChannel1 <- retDetail
// 			if delCount != len(deleteList) {
// 				logger.Errorf("eval delCount have err %d %d", delCount, len(deleteList))
// 			}
// 		} else {
// 			logger.Errorf("RemoveTokens have err %s", err.Error())
// 		}
// 	}
// 	logger.Infof("result Count timer %d %v", count, time.Now().UnixNano()/1e6)
// }

func (m *defaultMgr) Eval2(req *match_evaluator.ToEvalReq, keyy string, evalStartTime int64, resultChan chan *detailResult, inTeam *map[string]bool) {
	missInts := map[int]bool{}
	for index, detail := range req.Details {
		have := false
		for _, ply := range detail.Ids {
			_, ok := (*inTeam)[ply]
			have = have || ok
		}

		if !have {
			for _, ply := range detail.Ids {
				(*inTeam)[ply] = true
			}
		} else {
			missInts[index] = true
		}
	}
	nowV, err := m.getPoolVersion(keyy)
	if err == nil && nowV == req.Version {
		resultChan <- &detailResult{
			List:               req.Details,
			MissIndex:          missInts,
			EvalGroupSubId:     req.EvalGroupSubId,
			EvalGroupTaskCount: req.EvalGroupTaskCount,
			Version:            req.Version,
			Key:                keyy,
			StartTime:          req.StartTime,
			GameId:             req.GameId,
			SubType:            req.SubType,
			EvalStartTime:      evalStartTime,
			RunTime:            req.RunTime,
		}
	} else {
		logger.Infof("result %d poolversion miss", req.EvalGroupSubId)
	}
	//logger.Infof("result Count timer %d %v", req.EvalGroupSubId, time.Now().UnixNano()/1e6)
}

func (m *defaultMgr) Rem2(req *detailResult) {
	nowV, err := m.getPoolVersion(req.Key)
	if err != nil || nowV != req.Version {
		logger.Errorf("Rem2 poolversion 1 miss %d %d", nowV, req.Version)
		return
	}
	retDetail := make([]*match_evaluator.MatchDetail, 0, 256)
	needCount := 0

	innerFunc := func() bool {
		nowV, err := m.getPoolVersion(req.Key)
		if err == nil && nowV == req.Version {
			delCount, err := db.Default.RemoveTokens(context.Background(), retDetail, needCount, req.Key)
			if err == nil {
				m.resultChannel1 <- retDetail
				if delCount != needCount {
					logger.Errorf("eval delCount have err %d %d %d", req.EvalGroupSubId, delCount, needCount)
				} else {
					//logger.Infof("RemoveTokens success %d %d %d", req.EvalGroupSubId, len(retDetail), time.Now().UnixNano()/1e6)
				}
			} else {
				logger.Errorf("RemoveTokens have err %d %s", req.EvalGroupSubId, err.Error())
			}
			needCount = 0
			retDetail = retDetail[:0]
		} else {
			logger.Errorf("Rem2 poolversion 2 miss %d %d", nowV, req.Version)
			return false
		}
		return true
	}

	for pos, detail := range req.List {
		if _, ok := req.MissIndex[pos]; ok {
			continue
		}
		retDetail = append(retDetail, detail)
		needCount += len(detail.Ids)
		if needCount >= 512 {
			if !innerFunc() {
				return
			}
		}
	}
	if needCount > 0 {
		innerFunc()
	}
	//logger.Infof("RemoveTokens success %d %d %d", req.EvalGroupSubId, len(req.List), time.Now().UnixNano()/1e6)
}

func (m *defaultMgr) getPoolVersion(key string) (int64, error) {
	var err error
	var newV int64
	version, ok := m.versionMap.Load(key)
	if !ok {
		newV, err = db.Default.GetPoolVersion(context.Background(), key)
		if err != nil {
			logger.Errorf("get version from redis have err %s", err.Error())
			return 0, err
		} else {
			version = newV
			m.versionMap.LoadOrStore(key, version)
		}
	}
	return version.(int64), nil
}

func (m *defaultMgr) publishResult(detail *match_evaluator.MatchDetail) error {
	by, _ := proto.Marshal(detail)
	return broker.Publish("match_result", &broker.Message{
		Header: map[string]string{},
		Body:   by,
	})
}
