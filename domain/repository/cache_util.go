package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/smart-fm/kf-api/infrastructure/redis"
)

func getCacheByKey[T any](ctx context.Context, key string) *T {
	res, err := redis.GetRedisClient().Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	var entity T
	err = json.Unmarshal([]byte(res), &entity)
	if err != nil {
		return nil
	}
	return &entity
}

func deleteCacheByKey(ctx context.Context, key string) {
	redis.GetRedisClient().Del(ctx, key)
}

func setCacheByKey[T any](ctx context.Context, key string, entity T, expire time.Duration) {
	data, err := json.Marshal(entity)
	if err != nil {
		return
	}
	redis.GetRedisClient().Set(ctx, key, data, expire)
}
