package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/smart-fm/kf-api/config"
	"sync"
	"time"
)

var (
	redisClient     *redis.Client
	redisClientOnce sync.Once
)

func InitRedis() {
	redisClientOnce.Do(func() {
		conf := config.GetConfig()
		redisClient = redis.NewClient(&redis.Options{
			Network:        "tcp",
			Addr:           conf.RedisConfig.Address,
			ClientName:     "kf-api",
			Password:       conf.RedisConfig.Password,
			DB:             conf.RedisConfig.DB,
			MaxRetries:     3,
			DialTimeout:    60 * time.Second,
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			PoolFIFO:       true,
			PoolSize:       10,
			PoolTimeout:    60 * time.Second,
			MinIdleConns:   2,
			MaxIdleConns:   5,
			MaxActiveConns: 20,
		})

		res := redisClient.Ping(context.Background())
		if res.Err() != nil {
			panic("redis connect failed: " + res.Err().Error())
		}
	})
}

func GetRedisClient() *redis.Client {
	return redisClient
}
