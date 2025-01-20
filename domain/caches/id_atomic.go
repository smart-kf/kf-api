package caches

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/smart-fm/kf-api/infrastructure/redis"
)

const (
	idKey = "atomic.id"
)

type idAtomicCache struct{}

// GetOneId 全局自增id. 不管什么业务都可以直接用于主键.
func (c *idAtomicCache) GetOneId(ctx context.Context) (int64, error) {
	res, err := c.GetManyId(ctx, 1)
	if err != nil {
		return 0, err
	}
	if len(res) > 0 {
		return res[0], err
	}
	return 0, errors.New("getOneId failed")
}

// GetManyId 获取多个全局唯一id
func (c *idAtomicCache) GetManyId(ctx context.Context, n int) ([]int64, error) {
	cli := redis.GetRedisClient()
	result, err := cli.IncrBy(ctx, idKey, int64(n)).Result()
	if err != nil {
		return nil, err
	}
	var res = make([]int64, 0, n)
	for i := result - int64(n); i <= result; i++ {
		res = append(res, i)
	}
	return res, nil
}

// GetBizId 获取唯一主键id, 跟日期相关
// yyyymmddhhiiss[rand4] globalId
func (c *idAtomicCache) GetBizId(ctx context.Context) (int64, error) {
	gid, err := c.GetOneId(ctx)
	if err != nil {
		return 0, err
	}
	n := time.Now().Format(`60102150405`)
	randNum := rand.Intn(9000) + 1000

	x, _ := strconv.Atoi(fmt.Sprintf("%d%s%d", randNum, n, gid))
	return int64(x), nil
}

func (c *idAtomicCache) GetBizIds(ctx context.Context, n int) ([]int64, error) {
	gids, err := c.GetManyId(ctx, n)
	if err != nil {
		return nil, err
	}
	for index, val := range gids {
		n := time.Now().Format(`60102150405`)
		randNum := rand.Intn(9000) + 1000
		x, _ := strconv.Atoi(fmt.Sprintf("%d%s%d", randNum, n, val))
		gids[index] = int64(x)
	}
	return gids, nil
}
