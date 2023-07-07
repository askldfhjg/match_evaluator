package redis

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

const (
	allTickets     = "allTickets:%d:%s"
	poolVersionKey = "poolVersionKey:"
	taskFlag       = "taskFlag:%d:%s"
	poolLock       = "poolLock:%s"
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
	vv, err := redis.Int64(redisConn.Do("HGET", poolVersionKey+key, "poolVersionKey"))
	if err != nil {
		if err == redis.ErrNil {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return vv, nil
}

func (m *redisBackend) DelTaskFlag(ctx context.Context, key string, version int64, subTaskKey string) (int, error) {
	redisConn, err := m.redisPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer handleConnectionClose(&redisConn)

	script := `
	if(#ARGV < 1) then
		error("argv error")
	end
	local values = redis.call('HDEL', KEYS[1], ARGV[1])
	return redis.call('HLEN', KEYS[1])
	`

	args := []interface{}{subTaskKey}
	keys := []interface{}{fmt.Sprintf(taskFlag, version, key)}
	params := []interface{}{script, len(keys)}
	params = append(params, keys...)
	params = append(params, args...)
	return redis.Int(redisConn.Do("EVAL", params...))
}

func (m *redisBackend) UnLockPool(ctx context.Context, key string, version int64) error {
	redisConn, err := m.redisPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer handleConnectionClose(&redisConn)

	script := `
	if(#ARGV < 1) then
		error("argv error")
	end
	local value = redis.call('GET', KEYS[1])
	if(not value) then
		error("delete error")
	end
	if(value ~= ARGV[1]) then
		error("delete wrong version")
	end
	redis.call('DEL', KEYS[1])
	`
	args := []interface{}{version}
	keys := []interface{}{fmt.Sprintf(poolLock, key)}
	params := []interface{}{script, len(keys)}
	params = append(params, keys...)
	params = append(params, args...)
	_, err = redisConn.Do("EVAL", params...)
	return err
}
