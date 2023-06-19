package handler

import (
	"context"

	log "github.com/micro/micro/v3/service/logger"

	match_evaluator "match_evaluator/proto"
)

type Match_evaluator struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Match_evaluator) Call(ctx context.Context, req *match_evaluator.Request, rsp *match_evaluator.Response) error {
	log.Info("Received Match_evaluator.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}
