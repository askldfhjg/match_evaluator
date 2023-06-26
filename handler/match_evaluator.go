package handler

import (
	"context"

	"match_evaluator/evaluator"
	match_evaluator "match_evaluator/proto"
)

type Match_evaluator struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Match_evaluator) ToEval(ctx context.Context, req *match_evaluator.ToEvalReq, rsp *match_evaluator.ToEvalRsp) error {
	//log.Info("Received ToEval request")
	evaluator.DefaultManager.AddMsg(req)
	return nil
}
