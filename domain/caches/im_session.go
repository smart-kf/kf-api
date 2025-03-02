package caches

import (
	"context"
	"fmt"

	"github.com/smart-fm/kf-api/domain/factory"
	"github.com/smart-fm/kf-api/infrastructure/redis"
)

type imSessionCache struct{}

/*--------------------- 客服后台-------*/
// hashmap
func (imSessionCache) getKfBeOnlineSessionKey(cardId string) string {
	return fmt.Sprintf("kfbe_online_session_%s", cardId)
}

// SetKfbeOnlineSession 设置后台用户上线session  // kfbe_online_session_TM-oyTZ3v1toG
func (c *imSessionCache) SetKfbeOnlineSession(ctx context.Context, cardId string, sessionId string) {
	key := c.getKfBeOnlineSessionKey(cardId)
	cli := redis.GetRedisClient()
	cli.LPush(ctx, key, sessionId)
}

func (c *imSessionCache) DeleteKfBeOnlineSession(ctx context.Context, cardId string, sessionId string) {
	key := c.getKfBeOnlineSessionKey(cardId)
	cli := redis.GetRedisClient()
	cli.LRem(ctx, key, 1, sessionId)
}

// GetKfbeSessionIds 获取cardId对应的所有后端的sessionId
func (c *imSessionCache) GetKfbeSessionIds(ctx context.Context, cardId string) ([]string, error) {
	key := c.getKfBeOnlineSessionKey(cardId)
	cli := redis.GetRedisClient()
	val, err := cli.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

/*--------------------- 客服前台-------
前台用户的session 是一对一的, 只关心最新的session.
*/

func (imSessionCache) getKffeOnlineSessionKey(cardId string) string {
	return fmt.Sprintf("kffe_online_session_%s", cardId)
}

// SetKfbeOnlineSession 设置后台用户上线session
// 数据解构为 hashmap
// bucket = prefix.cardId
// bucket.key = uuid
// bucket.value = sessionId
func (c *imSessionCache) SetKffeOnlineSession(ctx context.Context, cardId string, userId string, sessionId string) {
	key := c.getKffeOnlineSessionKey(cardId)
	cli := redis.GetRedisClient()
	cli.HSet(ctx, key, userId, sessionId)
}

func (c *imSessionCache) DeleteKffeOnlineSession(ctx context.Context, cardId string, userId string) {
	key := c.getKffeOnlineSessionKey(cardId)
	cli := redis.GetRedisClient()
	cli.HDel(ctx, key, userId)
}

// GetKffeSessionIds 获取前台所有 对应 uids 用户的 sessionIds.
// 如果uids 为空，则获取所有的.
func (c *imSessionCache) GetKffeSessionIds(ctx context.Context, cardId string, uids []string) (
	map[string]string,
	error,
) {
	key := c.getKffeOnlineSessionKey(cardId)
	cli := redis.GetRedisClient()
	var (
		res map[string]string
		err error
	)
	if len(uids) == 0 {
		res, err = cli.HGetAll(ctx, key).Result()
	} else {
		var val []interface{}
		val, err = cli.HMGet(ctx, key, uids...).Result()
		if err != nil {
			return nil, err
		}
		res = make(map[string]string)
		for i, value := range val {
			if value == nil {
				res[uids[i]] = ""
			} else {
				res[uids[i]] = value.(string)
			}
		}
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

/* 下面是一些组合方法，对前后台的方法进行聚合utils */

// GetCardIDByKFFEToken 通过前台token获取cardId
func (c *imSessionCache) GetCardIDByKFFEToken(ctx context.Context, token string) (string, error) {
	// 1. 通过token 拿到 kf_card.id
	cardMainId, err := factory.FactoryParseUserToken(token)
	if err != nil {
		return "", err
	}
	// 2. 需要通过 kf_card.id 从redis中获取cardId
	cardId, err := KfCardCacheInstance.GetCardIDByMainID(ctx, cardMainId)
	if err != nil {
		return "", err
	}
	return cardId, nil
}

// GetKFFESessionIdByUserId 通过userId、cardId拿到前台用户的sessionId
// ps:: 可能为空哦.
func (c *imSessionCache) GetKFFESessionIdByUserId(ctx context.Context, cardId string, userId string) (string, error) {
	result, err := c.GetKffeSessionIds(ctx, cardId, []string{userId})
	if err != nil {
		return "", err
	}
	if len(result) > 0 {
		return result[userId], nil
	}
	return "", nil
}
