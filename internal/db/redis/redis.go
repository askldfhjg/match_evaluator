package redis

import (
	"context"
	match_evaluator "match_evaluator/proto"

	"github.com/gomodule/redigo/redis"
)

const (
	allTickets = "allTickets:"
	//ticketKey      = "ticket:%s"
	poolVersionKey = "poolVersionKey:"
)

func (m *redisBackend) RemoveTokens(ctx context.Context, retDetail []*match_evaluator.MatchDetail, needCount int, key string) (int, error) {
	redisConn, err := m.redisPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer handleConnectionClose(&redisConn)
	zsetKey := allTickets + key

	queryParams := make([]interface{}, needCount+1)
	index := 0
	queryParams[index] = zsetKey

	for _, detail := range retDetail {
		for _, id := range detail.Ids {
			if id == "robot" {
				continue
			}
			index++
			queryParams[index] = id
		}
	}
	return redis.Int(redisConn.Do("ZREM", queryParams...))
}

func (m *redisBackend) GetPoolVersion(ctx context.Context, key string) (int64, error) {
	redisConn, err := m.redisPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer handleConnectionClose(&redisConn)
	vv, err := redis.Int64(redisConn.Do("GET", poolVersionKey+key))
	if err != nil {
		if err == redis.ErrNil {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return vv, nil
}
