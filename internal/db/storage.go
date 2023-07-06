package db

import (
	"context"
)

var Default Service

type Service interface {
	Init(ctx context.Context, opts ...Option) error
	Close(ctx context.Context) error
	String() string
	//RemoveTokens(ctx context.Context, retDetail []*match_evaluator.MatchDetail, needCount int, key string) (int, error)
	GetPoolVersion(ctx context.Context, key string) (int64, error)
	MoveTokens(ctx context.Context, version int64, retDetail map[string]int32, key string) (int, error)
}

type MatchInfo struct {
	Id       string
	PlayerId string
	Score    int64
	GameId   string
	SubType  int64
}
