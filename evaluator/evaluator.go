package evaluator

import match_evaluator "match_evaluator/proto"

var DefaultManager Manager

type Manager interface {
	Start() error
	Stop() error
	AddMsg(req *match_evaluator.ToEvalReq)
}
