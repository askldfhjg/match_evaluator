package db

import (
	"context"
)

var Default Service

type Service interface {
	Init(ctx context.Context, opts ...Option) error
	Close(ctx context.Context) error
	String() string
	RemoveTokens(ctx context.Context, playerIds []string, gameId string, subType int64) (int, error)
	GetPoolVersion(ctx context.Context, key string) (int64, error)
}

type MatchInfo struct {
	Id       string
	PlayerId string
	Score    int64
	GameId   string
	SubType  int64
}
