package db

import (
	"context"
	match_evaluator "match_evaluator/proto"
)

var Default Service

type Service interface {
	Init(ctx context.Context, opts ...Option) error
	Close(ctx context.Context) error
	String() string
	RemoveTokens(ctx context.Context, retDetail []*match_evaluator.MatchDetail, needCount int, key string) (int, error)
	GetPoolVersion(ctx context.Context, key string) (int64, error)
}

type MatchInfo struct {
	Id       string
	PlayerId string
	Score    int64
	GameId   string
	SubType  int64
}
