package caches

import (
	"context"
	"fmt"
	"time"

	"github.com/smart-fm/kf-api/infrastructure/redis"
)

type kfAuthCache struct{}

func (c *kfAuthCache) getKey(token string) string {
	return fmt.Sprintf("kfbe_auth.%s", token)
}

func (c *kfAuthCache) getFrontKey(token string) string {
	return fmt.Sprintf("kfbe_fe_auth.%s", token)
}

func (c *kfAuthCache) SetBackendToken(ctx context.Context, token string, cardId string) error {
	redisClient := redis.GetRedisClient()
	err := redisClient.Set(ctx, c.getKey(token), cardId, time.Hour*24).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *kfAuthCache) GetBackendToken(ctx context.Context, token string) (string, error) {
	redisClient := redis.GetRedisClient()
	res, err := redisClient.Get(ctx, c.getKey(token)).Result()
	if err != nil {
		return "", fmt.Errorf("token not found")
	}
	return res, nil
}

func (c *kfAuthCache) SetFrontToken(ctx context.Context, token string, cardId string) error {
	redisClient := redis.GetRedisClient()
	err := redisClient.Set(ctx, c.getFrontKey(token), cardId, time.Hour*24).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *kfAuthCache) GetFrontToken(ctx context.Context, token string) (string, error) {
	redisClient := redis.GetRedisClient()
	res, _ := redisClient.Get(ctx, c.getFrontKey(token)).Result()
	if res == "" {
		return "", fmt.Errorf("token not found")
	}
	return res, nil
}
