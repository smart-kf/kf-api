package caches

import (
	"context"
	"encoding/json"
	"time"

	redis2 "github.com/redis/go-redis/v9"

	"github.com/smart-fm/kf-api/infrastructure/redis"
)

func UnMarshalRedisResult[T any](cmd *redis2.StringCmd, t *T, missFunc func() (*T, error)) error {
	data, err := cmd.Bytes()
	if err != nil {
		if err == redis2.Nil {
			if missFunc != nil {
				newT, err := missFunc()
				if err != nil {
					return err
				}
				*t = *newT
			}
			return nil
		}
		return err
	}
	return json.Unmarshal(data, t)
}

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
