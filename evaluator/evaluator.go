package evaluator

var DefaultManager Manager

type Manager interface {
	Start() error
	Stop() error
	AddMsg(req interface{})
}
