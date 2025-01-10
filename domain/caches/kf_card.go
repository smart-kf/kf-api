package caches

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/infrastructure/redis"
)

const kfCardCacheKey = "kf_card_cache_%s"
const kfCardMainIdCacheKey = "kf_card_main_id_%d"

type kfCardCacheInstance struct{}

func (c *kfCardCacheInstance) getKey(cardID string) string {
	return fmt.Sprintf(kfCardCacheKey, cardID)
}

func (c *kfCardCacheInstance) getMainIdKey(id int64) string {
	return fmt.Sprintf(kfCardMainIdCacheKey, id)
}

// Delete 删除缓存.
func (c *kfCardCacheInstance) Delete(ctx context.Context, cardId string) error {
	deleteCacheByKey(ctx, c.getKey(cardId))
	return nil
}

// GetCardByID 获取卡密缓存，可能击穿到db
func (c *kfCardCacheInstance) GetCardByID(ctx context.Context, cardId string) (*dao.KFCard, error) {
	cli := redis.GetRedisClient()
	cmd := cli.Get(ctx, c.getKey(cardId))
	var res dao.KFCard
	err := UnMarshalRedisResult[dao.KFCard](
		cmd, &res, func() (*dao.KFCard, error) {
			cardRepo := repository.KFCardRepository{}
			kfCardDao, ok, err := cardRepo.FindByCardID(ctx, cardId)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, fmt.Errorf("cardId not found: %+v", cardId)
			}
			_ = c.SetCardCacheByCardID(ctx, kfCardDao)
			return kfCardDao, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// SetCardCacheByCardID 设置缓存，失效时间为：[30,60) min
func (c *kfCardCacheInstance) SetCardCacheByCardID(ctx context.Context, card *dao.KFCard) error {
	key := c.getKey(card.CardID)
	setCacheByKey(ctx, key, card, time.Duration(rand.Intn(30)+30)*time.Minute)
	return nil
}

// GetCardIDByMainID 这里给前台用户，做一个从 kf_card.id -> cardId 的映射, 缓存方式, 击穿到数据库
// 前台用户的token规则为：mainId|uuid, mainId 暴露在前端.
func (c *kfCardCacheInstance) GetCardIDByMainID(ctx context.Context, id int64) (string, error) {
	cacheKey := c.getMainIdKey(id)
	str := getCacheByKey[string](ctx, cacheKey)
	if str != nil {
		if *str != "" {
			return *str, nil
		}
	}
	var repo repository.KFCardRepository
	card, ok, err := repo.FindByMainId(ctx, id)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("card not found:%d", id)
	}
	setCacheByKey[string](ctx, cacheKey, card.CardID, 24*time.Hour)
	return card.CardID, nil
}
