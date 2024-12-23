package repository

import (
	"context"
	"encoding/json"
	"std-api/pkg/db"
	"time"
)

func getCacheByKey[T any](ctx context.Context, key string) *T {
	res, err := db.GetRedisClient().Get(ctx, key).Result()
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
	db.GetRedisClient().Del(ctx, key)
}

func setCacheByKey[T any](ctx context.Context, key string, entity T, expire time.Duration) {
	data, err := json.Marshal(entity)
	if err != nil {
		return
	}
	db.GetRedisClient().Set(ctx, key, data, expire)
}
