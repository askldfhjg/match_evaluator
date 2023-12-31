package handler

import (
	"context"

	"match_evaluator/evaluator"
	match_evaluator "match_evaluator/proto"
)

type Match_evaluator struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Match_evaluator) ToEval(ctx context.Context, req *match_evaluator.ToEvalReq, rsp *match_evaluator.ToEvalRsp) error {
	//logger.Infof("ToEval timer %v", time.Now().UnixNano()/1e6)
	evaluator.DefaultManager.AddMsg(req)
	return nil
}

func (e *Match_evaluator) ToEvalReady(ctx context.Context, req *match_evaluator.ToEvalReadyReq, rsp *match_evaluator.ToEvalReadyRsp) error {
	//logger.Infof("ToEvalReady timer %v", time.Now().UnixNano()/1e6)
	evaluator.DefaultManager.AddMsg(req)
	return nil
}
