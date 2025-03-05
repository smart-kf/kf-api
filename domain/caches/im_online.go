package caches

import (
	"context"
	"fmt"

	"github.com/smart-fm/kf-api/infrastructure/redis"
)

/* 数据解构为 set */
/*
	bucket = kf_online_$cardId
	sets = (userId1,userId2,userId3...)
*/

const (
	onLineKey = "kf_online_%s"
)

type userOnLineCache struct{}

func (i *userOnLineCache) getKey(cardId string) string {
	key := fmt.Sprintf(onLineKey, cardId)
	return key
}

// SetUserOnline 设置用户在线.
func (d *userOnLineCache) SetUserOnline(ctx context.Context, cardId string, userId string) error {
	_, err := redis.GetRedisClient().SAdd(ctx, d.getKey(cardId), userId).Result()
	return err
}

// SetUserOnline 设置用户离线
func (d *userOnLineCache) SetUserOffline(ctx context.Context, cardId string, userId string) error {
	_, err := redis.GetRedisClient().SRem(ctx, d.getKey(cardId), userId).Result()
	return err
}

// IsUserOnline 判断userid 是否在线.
func (d *userOnLineCache) IsUserOnline(ctx context.Context, cardId string, userId string) (bool, error) {
	key := d.getKey(cardId)
	ok, err := redis.GetRedisClient().SIsMember(ctx, key, userId).Result()
	return ok, err
}

// IsUserOnline 批量判断user在线情况
func (d *userOnLineCache) IsUsersOnline(ctx context.Context, cardId string, userId []string) (map[string]bool, error) {
	if len(userId) == 0 {
		return make(map[string]bool), nil
	}
	script := `local results = {}
for i = 1, #KEYS do
    results[i] = redis.call("SISMEMBER", ARGV[1], KEYS[i])
end
return results`
	key := d.getKey(cardId)
	args := make([]interface{}, len(userId)+1)
	args[0] = key
	for i, id := range userId {
		args[i+1] = id
	}
	iface, err := redis.GetRedisClient().Eval(ctx, script, userId, args...).Result()
	if err != nil {
		return nil, err
	}
	res := make(map[string]bool)
	if resultArray, ok := iface.([]interface{}); ok {
		for i, exists := range resultArray {
			if exists.(int64) == 1 {
				res[userId[i]] = true
			} else {
				res[userId[i]] = false
			}
		}
	}
	return res, nil
}

// GetOnLineUsers 获取所有在线用户的 uuid
func (d *userOnLineCache) GetOnLineUsers(ctx context.Context, cardId string) ([]string, error) {
	key := d.getKey(cardId)
	res, err := redis.GetRedisClient().SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}
