package redis

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func (m *redisBackend) RemoveTokens(ctx context.Context, playerIds []string, gameId string, subType int64) (int, error) {
	redisConn, err := m.redisPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer handleConnectionClose(&redisConn)
	zsetKey := fmt.Sprintf(allTickets, gameId, subType)
	//inter1 := make([]interface{}, 0, len(playerIds))
	inter2 := make([]interface{}, 0, len(playerIds))
	inter2 = append(inter2, zsetKey)
	for _, ply := range playerIds {
		//inter1 = append(inter1, fmt.Sprintf(ticketKey, ply))
		inter2 = append(inter2, ply)
	}
	//delCount, _ := redis.Int(redisConn.Do("DEL", inter1...))
	return redis.Int(redisConn.Do("ZREM", inter2...))
}

func (m *redisBackend) GetPoolVersion(ctx context.Context, gameId string, subType int64) (int64, error) {
	redisConn, err := m.redisPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer handleConnectionClose(&redisConn)
	vv, err := redis.Int64(redisConn.Do("GET", fmt.Sprintf(poolVersionKey, gameId, subType)))
	if err != nil {
		if err == redis.ErrNil {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return vv, nil
}
