package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type ZSetMember struct {
	Member string
	Score  int64
}

func ZAdd(ctx context.Context, client *redis.Client, key string, members ...ZSetMember) error {
	zs := make([]redis.Z, 0, len(members))
	for _, v := range members {
		z := redis.Z{
			Score:  float64(v.Score),
			Member: v.Member,
		}
		zs = append(zs, z)
	}
	return client.ZAdd(ctx, key, zs...).Err()
}

// ZRangeByScore 获取 从 [now-step, now] 时间范围内的数据.
func ZRangeByScore(ctx context.Context, client *redis.Client, key string, start int64, step int64) ([]string, error) {
	res := client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", start-step),
		Max: fmt.Sprintf("%d", start),
	})
	return res.Result()
}
