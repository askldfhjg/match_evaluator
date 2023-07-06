package redis

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

const (
	allTickets     = "allTickets:%d:%s"
	poolVersionKey = "poolVersionKey:"
)

// func (m *redisBackend) RemoveTokens(ctx context.Context, retDetail []*match_evaluator.MatchDetail, needCount int, key string) (int, error) {
// 	redisConn, err := m.redisPool.GetContext(ctx)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer handleConnectionClose(&redisConn)
// 	zsetKey := allTickets + key

// 	queryParams := make([]interface{}, needCount+1)
// 	index := 0
// 	queryParams[index] = zsetKey

// 	for _, detail := range retDetail {
// 		for _, id := range detail.Ids {
// 			if id == "robot" {
// 				continue
// 			}
// 			index++
// 			queryParams[index] = id
// 		}
// 	}
// 	return redis.Int(redisConn.Do("ZREM", queryParams...))
// }

func (m *redisBackend) MoveTokens(ctx context.Context, version int64, retDetail map[string]int32, key string) (int, error) {
	redisConn, err := m.redisPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer handleConnectionClose(&redisConn)
	if len(retDetail) <= 0 {
		return 0, err
	}
	zsetKey := fmt.Sprintf(allTickets, version, key)

	queryParams := make([]interface{}, len(retDetail)*2+1)
	index := 0
	queryParams[index] = zsetKey

	for id, score := range retDetail {
		index++
		queryParams[index] = score
		index++
		queryParams[index] = id
	}
	return redis.Int(redisConn.Do("ZADD", queryParams...))
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
