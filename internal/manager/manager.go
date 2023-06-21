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

const processExpireTime = 10

func NewManager(opts ...evaluator.EvaluatorOption) evaluator.Manager {
	m := &defaultMgr{
		exited:      make(chan struct{}, 1),
		msgList:     make(chan *match_evaluator.ToEvalReq, 1000),
		msgBackList: make(chan string, 100),
		channelMap:  make(map[string]chan *match_evaluator.ToEvalReq, 100),
		versionMap:  cmap.New[int64](),
	}
	for _, o := range opts {
		o(&m.opts)
	}
	return m
}

type defaultMgr struct {
	opts        evaluator.EvaluatorOptions
	exited      chan struct{}
	msgList     chan *match_evaluator.ToEvalReq
	msgBackList chan string
	channelMap  map[string]chan *match_evaluator.ToEvalReq
	versionMap  cmap.ConcurrentMap[string, int64]
}

func (m *defaultMgr) Start() error {
	err := m.SubPoolVersion()
	if err != nil {
		return err
	}
	go m.loop()
	return nil
}

func (m *defaultMgr) Stop() error {
	close(m.exited)
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
		case <-m.exited:
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

func (m *defaultMgr) AddMsg(req *match_evaluator.ToEvalReq) {
	m.msgList <- req
}

func (m *defaultMgr) processEval(chnl chan *match_evaluator.ToEvalReq, gId string) {
	getIndex := 0
	groupId := gId
	channel := chnl
	var list []*match_evaluator.MatchDetail
	timer := time.After(10 * time.Second)
	go func() {
		for {
			select {
			case <-timer:
				m.msgBackList <- groupId
				list = nil
				return
			case req := <-channel:
				flag := (getIndex >> (req.EvalGroupSubId - 1)) & 1
				if flag <= 0 {
					getIndex += int(math.Exp2(float64(req.EvalGroupSubId - 1)))
					list = append(list, req.Details...)
				}

				if getIndex == int(math.Exp2(float64(req.EvalGroupTaskCount))) {
					m.eval(list, groupId, req.Version, req.GetGameId(), req.GetSubType())
					m.msgBackList <- groupId
					list = nil
					return
				}
			}
		}
	}()
}

func (m *defaultMgr) eval(list []*match_evaluator.MatchDetail, groupId string, version int64, gameId string, subType int64) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Score > list[j].Score
	})

	inTeam := map[string]int{}
	ctx := context.Background()
	for _, detail := range list {
		have := false
		for _, ply := range detail.Ids {
			_, ok := inTeam[ply]
			have = have || ok
		}
		if !have {
			for _, ply := range detail.Ids {
				inTeam[ply] = 1
			}
			nowV, err := m.getPoolVersion(gameId, subType)
			if err == nil && nowV == version {
				err = db.Default.RemoveTokens(ctx, detail.Ids, gameId, subType)
				if err == nil {
					errs := m.PublishPoolVersion(detail)
					if errs != nil {
						logger.Errorf("PublishPoolVersion have err %s", errs.Error())
					}
				} else {
					logger.Errorf("RemoveTokens have err %s", err.Error())
				}
			}
		}
	}
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

func (m *defaultMgr) PublishPoolVersion(detail *match_evaluator.MatchDetail) error {
	by, _ := proto.Marshal(detail)
	return broker.Publish("pool_version", &broker.Message{
		Header: map[string]string{},
		Body:   by,
	})
}