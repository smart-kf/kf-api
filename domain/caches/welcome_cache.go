package caches

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/infrastructure/redis"
)

var (
	cacheKey              = "kf:welcome:message:%s:%s"      // 存储后台的欢迎语缓存
	needSendWelcomeMsgKey = "kf:user:welcome:message:%s:%s" // 存储新用户的key，过期时间10s钟，代表这个客户需要发送欢迎语
)

type WelcomeMessageCache struct{}

func (c *WelcomeMessageCache) getKey(cardId string, msgType string) string {
	return fmt.Sprintf(cacheKey, cardId, msgType)
}

func (c *WelcomeMessageCache) getNeedSendWelcomeMsgKey(cardID string, userID string) string {
	return fmt.Sprintf(needSendWelcomeMsgKey, cardID, userID)
}

func (c *WelcomeMessageCache) SetUserNeedSendMsg(ctx context.Context, cardId string, userId string) {
	key := c.getNeedSendWelcomeMsgKey(cardId, userId)
	redis.GetRedisClient().Set(ctx, key, 1, time.Second*10)
}

func (c *WelcomeMessageCache) GetUserNeedSendMsg(ctx context.Context, cardId string, userId string) bool {
	key := c.getNeedSendWelcomeMsgKey(cardId, userId)
	_, err := redis.GetRedisClient().Get(ctx, key).Result()
	if err != nil {
		return false
	}
	return true
}

func (c *WelcomeMessageCache) DelUserNeedSendMsg(ctx context.Context, cardId string, userId string) {
	key := c.getNeedSendWelcomeMsgKey(cardId, userId)
	redis.GetRedisClient().Del(ctx, key)
}

func (c *WelcomeMessageCache) FindWelcomeMessages(ctx context.Context, cardId string) []*dao.KfWelcomeMessage {
	key := c.getKey(cardId, constant.WelcomeMsg)
	result, err := redis.GetRedisClient().Get(ctx, key).Bytes()
	if err != nil {
		msgs := c.findFromRepository(ctx, cardId, constant.WelcomeMsg)
		if len(msgs) > 0 {
			data, _ := json.Marshal(msgs)
			redis.GetRedisClient().Set(ctx, key, string(data), time.Minute*5)
		}

		return msgs
	}
	var msgs []*dao.KfWelcomeMessage
	json.Unmarshal(result, &msgs)

	return msgs
}

func (c *WelcomeMessageCache) FindSmartMsg(
	ctx context.Context,
	cardId string,
	primaryId int64,
) *dao.KfWelcomeMessage {
	key := c.getKey(cardId, fmt.Sprintf("%s:%d", constant.SmartMsg, primaryId))
	result, err := redis.GetRedisClient().Get(ctx, key).Bytes()
	if err != nil {
		msg := c.findSmartMsgFromRepository(ctx, cardId, primaryId)
		if msg == nil {
			return nil
		}
		data, _ := json.Marshal(msg)
		redis.GetRedisClient().Set(ctx, key, string(data), 5*time.Minute)
		return msg
	}
	var msgs dao.KfWelcomeMessage
	json.Unmarshal(result, &msgs)
	return &msgs
}

func (c *WelcomeMessageCache) DeleteCache(ctx context.Context, cardId string) {
	redis.GetRedisClient().Del(context.Background(), c.getKey(cardId, constant.WelcomeMsg))
}

func (c *WelcomeMessageCache) findFromRepository(
	ctx context.Context, cardId string,
	typ string,
) []*dao.KfWelcomeMessage {
	var repo repository.KfWelcomeMessageRepository
	msg, _, err := repo.List(ctx, cardId, constant.WelcomeMsg, 1, 10)
	if err != nil {
		return nil
	}
	return msg
}

func (c *WelcomeMessageCache) findSmartMsgFromRepository(
	ctx context.Context, cardId string,
	id int64,
) *dao.KfWelcomeMessage {
	var repo repository.KfWelcomeMessageRepository
	msg, err := repo.FindById(ctx, cardId, id)
	if err != nil {
		return nil
	}
	return msg
}
