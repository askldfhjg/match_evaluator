package manager

import (
	"context"
	"fmt"
	"match_evaluator/evaluator"
	"match_evaluator/internal/db"
	match_evaluator "match_evaluator/proto"
	"math"
	"sort"
	"sync"
	"time"

	match_process "github.com/askldfhjg/match_apis/match_process/proto"

	"github.com/micro/micro/v3/service/broker"
	"github.com/micro/micro/v3/service/logger"
	"google.golang.org/protobuf/proto"
)

const deleteGroupCount = 500

func NewManager(opts ...evaluator.EvaluatorOption) evaluator.Manager {
	m := &defaultMgr{
		exited1:        make(chan struct{}, 1),
		exited2:        make(chan struct{}, 1),
		exited3:        make(chan struct{}, 1),
		msgList:        make(chan *match_evaluator.ToEvalReq, 1000),
		msgBackList:    make(chan string, 100),
		resultChannel1: make(chan []*match_evaluator.MatchDetail, 1000),
		resultChannel2: make(chan *match_evaluator.MatchDetail, 100000),
		channelMap:     make(map[string]chan *match_evaluator.ToEvalReq, 100),
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
	msgList        chan *match_evaluator.ToEvalReq
	msgBackList    chan string
	resultChannel1 chan []*match_evaluator.MatchDetail
	resultChannel2 chan *match_evaluator.MatchDetail
	channelMap     map[string]chan *match_evaluator.ToEvalReq
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
			key := fmt.Sprintf("%s:%d", req.GameId, req.SubType)
			version, err := m.getPoolVersion(key)
			if err != nil {
				logger.Errorf("get version from redis have err %s", err.Error())
			} else if version == req.Version {
				channel, ok := m.channelMap[req.EvalGroupId]
				if !ok {
					channel = make(chan *match_evaluator.ToEvalReq, 10)
					m.channelMap[req.EvalGroupId] = channel
					go m.processEval(channel, req.EvalGroupId)
				}
				channel <- req
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

func (m *defaultMgr) AddMsg(req *match_evaluator.ToEvalReq) {
	m.msgList <- req
}

func (m *defaultMgr) processEval(chnl chan *match_evaluator.ToEvalReq, gId string) {
	getIndex := 0
	groupId := gId
	channel := chnl
	gameId := ""
	var taskCount int64
	var subType int64
	var version int64
	var list []*match_evaluator.MatchDetail
	timer := time.After(2 * time.Second)
	go func() {
		for {
			select {
			case <-timer:
				if len(list) > 0 {
					m.Eval(list, groupId, version, gameId, subType, taskCount == 1)
				}
				m.msgBackList <- groupId
				list = nil
				return
			case req := <-channel:
				gameId = req.GetGameId()
				subType = req.GetSubType()
				version = req.GetVersion()
				taskCount = req.GetEvalGroupTaskCount()
				flag := (getIndex >> (req.EvalGroupSubId - 1)) & 1
				if flag <= 0 {
					//logger.Infof("add item %v", req.EvalGroupSubId)
					getIndex += int(math.Exp2(float64(req.EvalGroupSubId - 1)))
					list = append(list, req.Details...)
				}
				//logger.Infof("now item %v", getIndex)
				if getIndex == int(math.Exp2(float64(req.EvalGroupTaskCount))-1) {
					//logger.Infof("start %v", req.EvalGroupTaskCount)
					m.msgBackList <- groupId
					m.Eval(list, groupId, req.Version, gameId, subType, taskCount == 1)
					//logger.Infof("result timer %v", time.Now().UnixNano()/1e6)
					list = nil
					return
				}
			}
		}
	}()
}

func (m *defaultMgr) Eval(list []*match_evaluator.MatchDetail, groupId string, version int64, gameId string, subType int64, pass bool) {
	if !pass {
		sort.Slice(list, func(i, j int) bool {
			return list[i].Score < list[j].Score
		})
	}
	keyy := fmt.Sprintf("%s:%d", gameId, subType)
	inTeam := map[string]int{}
	ctx := context.Background()
	count := 0
	deleteList := make([]string, 0, deleteGroupCount)
	retDetail := make([]*match_evaluator.MatchDetail, 0, deleteGroupCount)
	for _, detail := range list {
		have := false
		if !pass {
			for _, ply := range detail.Ids {
				_, ok := inTeam[ply]
				have = have || ok
			}
		}

		if !have || pass {
			if !pass {
				for _, ply := range detail.Ids {
					inTeam[ply] = 1
				}
			}
			nowV, err := m.getPoolVersion(keyy)
			if err == nil && nowV == version {
				deleteList = append(deleteList, detail.Ids...)
				retDetail = append(retDetail, detail)
				if len(deleteList) >= deleteGroupCount {
					delCount, err := db.Default.RemoveTokens(ctx, deleteList, gameId, subType)
					if err == nil {
						count += len(retDetail)
						m.resultChannel1 <- retDetail
						if delCount != len(deleteList) {
							logger.Errorf("eval delCount have err %d %d", delCount, len(deleteList))
						}
					} else {
						logger.Errorf("RemoveTokens have err %s", err.Error())
						break
					}
					deleteList = make([]string, 0, deleteGroupCount)
					retDetail = make([]*match_evaluator.MatchDetail, 0, deleteGroupCount)
				}
			} else {
				logger.Infof("result Count version break")
				break
			}
		}
	}
	if len(deleteList) > 0 {
		delCount, err := db.Default.RemoveTokens(ctx, deleteList, gameId, subType)
		if err == nil {
			count += len(retDetail)
			m.resultChannel1 <- retDetail
			if delCount != len(deleteList) {
				logger.Errorf("eval delCount have err %d %d", delCount, len(deleteList))
			}
		} else {
			logger.Errorf("RemoveTokens have err %s", err.Error())
		}
	}
	logger.Infof("result Count timer %d %v", count, time.Now().UnixNano()/1e6)
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
