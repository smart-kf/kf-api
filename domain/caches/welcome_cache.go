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
	cacheKey = "kf:welcome:message:%s:%s"
)

type WelcomeMessageCache struct{}

func (c *WelcomeMessageCache) getKey(cardId string, msgType string) string {
	return fmt.Sprintf(cacheKey, cardId, msgType)
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
