package caches

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

// key->value
const (
	kfSettingCacheKey = "kf.setting.%s"
)

type kfSettingCache struct{}

func (c *kfSettingCache) GetOne(ctx context.Context, cardID string) (*dao.KFSettings, error) {
	res := getCacheByKey[dao.KFSettings](ctx, fmt.Sprintf(kfSettingCacheKey, cardID))
	if res != nil {
		return res, nil
	}
	var repo repository.KFSettingRepository
	setting, ok, err := repo.GetByCardID(ctx, cardID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("not found")
	}
	setCacheByKey[dao.KFSettings](ctx, fmt.Sprintf(kfSettingCacheKey, cardID), *setting, 10*time.Minute)
	return setting, nil
}

func (c *kfSettingCache) DeleteOne(ctx context.Context, cardId string) {
	deleteCacheByKey(ctx, fmt.Sprintf(kfSettingCacheKey, cardId))
}
