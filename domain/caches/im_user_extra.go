package caches

import (
	"context"
	"encoding/json"
	"fmt"

	redis2 "github.com/redis/go-redis/v9"

	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/infrastructure/redis"
)

/* 持久化一些用户的扩展属性，用redis存储 */
/* 数据结构为 hashMap
bucket = kf_user_extra_$cardId
bucket.key = $userId -> user的uuid
bucket.value = $dao.UserExtra
*/
const (
	kfUserExtraKey = "kf_user_extra_%s"
)

type kfUserExtraCache struct{}

func (d *kfUserExtraCache) getKey(cardId string) string {
	key := fmt.Sprintf(kfUserExtraKey, cardId)
	return key
}

func (d *kfUserExtraCache) GetUserObj(ctx context.Context, cardId string, userId string) (dao.UserExtra, error) {
	key := d.getKey(cardId)
	data, err := redis.GetRedisClient().HGet(ctx, key, userId).Bytes()
	if err != nil {
		if err == redis2.Nil {
			return dao.UserExtra{}, nil
		}
		return dao.UserExtra{}, err
	}
	var userObj dao.UserExtra
	err = json.Unmarshal(data, &userObj)
	if err != nil {
		return dao.UserExtra{}, err
	}
	return userObj, nil
}

func (d *kfUserExtraCache) GetUserObjs(ctx context.Context, cardId string, userIds []string) (
	map[string]dao.UserExtra,
	error,
) {
	key := d.getKey(cardId)
	data, err := redis.GetRedisClient().HMGet(ctx, key, userIds...).Result()
	if err != nil {
		return nil, err
	}
	var userObj = make(map[string]dao.UserExtra)
	for index, item := range data {
		if item == nil {
			continue
		}
		var extra dao.UserExtra
		err = json.Unmarshal([]byte(item.(string)), &extra)
		if err != nil {
			return nil, err
		}
		userObj[userIds[index]] = extra
	}
	return userObj, nil
}

func (d *kfUserExtraCache) SetUserObj(ctx context.Context, cardId string, userId string, obj dao.UserExtra) error {
	key := d.getKey(cardId)
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return redis.GetRedisClient().HSet(ctx, key, userId, string(data)).Err()
}
