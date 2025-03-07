package caches

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/smart-fm/kf-api/infrastructure/redis"
)

// 过期时间为2天.
// 数据结构为 hashMap
// bucket = kf_user_unread_$cardId
// bucket.key = $userId
// bucket.value = $unReadMsgCount
const (
	userUnReadKey = "kf_user_unread_%s"
)

type userUnReadCache struct{}

func (d *userUnReadCache) getKey(cardId string) string {
	return fmt.Sprintf(userUnReadKey, cardId)
}

// IncrUserUnRead 增加用户的未读消息数量
// 如果 为 -1， 则删除对应的key.
func (d *userUnReadCache) IncrUserUnRead(ctx context.Context, cardId string, userId string, n int64) error {
	key := d.getKey(cardId)
	if n == -1 {
		return redis.GetRedisClient().HDel(ctx, key, userId).Err()
	}
	redis.GetRedisClient().HIncrBy(ctx, key, userId, n).Err()

	// 设置过期时间
	redis.GetRedisClient().Expire(ctx, key, time.Hour*48) // 设置2天的过期时间
	return nil
}

// GetUserUnRead 获取某个用户的未读消息数量
func (d *userUnReadCache) GetUserUnRead(ctx context.Context, cardId string, userId string) (int64, error) {
	key := d.getKey(cardId)
	return redis.GetRedisClient().HGet(ctx, key, userId).Int64()
}

// GetUnReadUsers 获取所有有未读消息数量的用户, 返回: map[userId]unReadMsgCount
func (d *userUnReadCache) GetUnReadUsers(ctx context.Context, cardId string) (map[string]int64, error) {
	key := d.getKey(cardId)
	result, err := redis.GetRedisClient().HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var res = make(map[string]int64)
	for k, v := range result {
		i, _ := strconv.ParseInt(v, 10, 64)
		res[k] = i
	}
	return res, nil
}
