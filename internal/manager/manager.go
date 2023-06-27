package manager

import (
	"context"
	"fmt"
	"match_evaluator/evaluator"
	"match_evaluator/internal/db"
	match_evaluator "match_evaluator/proto"
	"math"
	"sort"
	"time"

	match_process "github.com/askldfhjg/match_apis/match_process/proto"

	"github.com/micro/micro/v3/service/broker"
	"github.com/micro/micro/v3/service/logger"
	cmap "github.com/orcaman/concurrent-map/v2"
	"google.golang.org/protobuf/proto"
)

func NewManager(opts ...evaluator.EvaluatorOption) evaluator.Manager {
	m := &defaultMgr{
		exited1:       make(chan struct{}, 1),
		exited2:       make(chan struct{}, 1),
		msgList:       make(chan *match_evaluator.ToEvalReq, 1000),
		msgBackList:   make(chan string, 100),
		resultChannel: make(chan *match_evaluator.MatchDetail, 1000),
		channelMap:    make(map[string]chan *match_evaluator.ToEvalReq, 100),
		versionMap:    cmap.New[int64](),
	}
	for _, o := range opts {
		o(&m.opts)
	}
	return m
}

type defaultMgr struct {
	opts          evaluator.EvaluatorOptions
	exited1       chan struct{}
	exited2       chan struct{}
	msgList       chan *match_evaluator.ToEvalReq
	msgBackList   chan string
	resultChannel chan *match_evaluator.MatchDetail
	channelMap    map[string]chan *match_evaluator.ToEvalReq
	versionMap    cmap.ConcurrentMap[string, int64]
}

func (m *defaultMgr) Start() error {
	err := m.SubPoolVersion()
	if err != nil {
		return err
	}
	go m.resultPublishLoop()
	go m.loop()
	return nil
}

func (m *defaultMgr) Stop() error {
	close(m.exited1)
	close(m.exited2)
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
	logger.Info("handlerMsg %+v", msg)
	m.versionMap.Set(fmt.Sprintf("%s:%d", msg.GameId, msg.SubType), msg.Version)
	return nil
}

func (m *defaultMgr) loop() {
	for {
		select {
		case <-m.exited1:
			return
		case req := <-m.msgList:
			version, err := m.getPoolVersion(req.GameId, req.SubType)
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

func (m *defaultMgr) resultPublishLoop() {
	for {
		select {
		case <-m.exited2:
			return
		case req := <-m.resultChannel:
			err := m.publishResult(req)
			if err != nil {
				logger.Errorf("PublishResult have err %s", err.Error())
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
					m.eval(list, groupId, version, gameId, subType, taskCount == 1)
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
					m.eval(list, groupId, req.Version, gameId, subType, taskCount == 1)
					m.msgBackList <- groupId
					list = nil
					return
				}
			}
		}
	}()
}

func (m *defaultMgr) eval(list []*match_evaluator.MatchDetail, groupId string, version int64, gameId string, subType int64, pass bool) {
	if !pass {
		sort.Slice(list, func(i, j int) bool {
			return list[i].Score < list[j].Score
		})
	}

	inTeam := map[string]int{}
	ctx := context.Background()
	count := 0
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
			nowV, err := m.getPoolVersion(gameId, subType)
			if err == nil && nowV == version {
				delCount, err := db.Default.RemoveTokens(ctx, detail.Ids, gameId, subType)
				if err == nil {
					if delCount != len(detail.Ids) {
						logger.Errorf("eval delCount have err %d %d", delCount, len(detail.Ids))
					} else {
						m.resultChannel <- detail
						count++
					}
				} else {
					logger.Errorf("RemoveTokens have err %s", err.Error())
				}
			}
		}
	}
	logger.Infof("result Count %d %v", count, time.Now().UnixNano()/1e6)
}

func (m *defaultMgr) getPoolVersion(gameId string, subType int64) (int64, error) {
	var err error
	var newV int64
	key := fmt.Sprintf("%s:%d", gameId, subType)
	version, ok := m.versionMap.Get(key)
	if !ok {
		newV, err = db.Default.GetPoolVersion(context.Background(), gameId, subType)
		if err != nil {
			logger.Errorf("get version from redis have err %s", err.Error())
			return 0, err
		} else {
			version = newV
			m.versionMap.SetIfAbsent(key, version)
		}
	}
	return version, nil
}

func (m *defaultMgr) publishResult(detail *match_evaluator.MatchDetail) error {
	by, _ := proto.Marshal(detail)
	return broker.Publish("match_result", &broker.Message{
		Header: map[string]string{},
		Body:   by,
	})
}
