package handler

import (
	"context"
	"time"

	"match_evaluator/evaluator"
	match_evaluator "match_evaluator/proto"

	"github.com/micro/micro/v3/service/logger"
)

type Match_evaluator struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Match_evaluator) ToEval(ctx context.Context, req *match_evaluator.ToEvalReq, rsp *match_evaluator.ToEvalRsp) error {
	logger.Infof("ToEval timer %v", time.Now().UnixNano()/1e6)
	evaluator.DefaultManager.AddMsg(req)
	return nil
}
