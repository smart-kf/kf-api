package caches

import (
	"context"
	"fmt"
	"time"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

/*
存储数据库的用户信息
expire = 1day
key = kf_user_cache_$cardId_$user.UUID
*/
const (
	kfUserCacheKey = "kf_user_cache_%s_%s"
)

type kfUserCache struct{}

func (d *kfUserCache) getKey(cardId string, userId string) string {
	return fmt.Sprintf(kfUserCacheKey, cardId, userId)
}

func (d *kfUserCache) SetDBUser(ctx context.Context, cardId string, user *dao.KfUser) {
	key := d.getKey(cardId, user.UUID)
	setCacheByKey(ctx, key, user, 24*time.Hour)
}

func (d *kfUserCache) DelDBUserCache(ctx context.Context, cardId string, user *dao.KfUser) {
	key := d.getKey(cardId, user.UUID)
	deleteCacheByKey(ctx, key)
}

// GetDBUser 从缓存中获取数据，如果未找到会查询数据库.
// 未做攻击防范.
func (d *kfUserCache) GetDBUser(ctx context.Context, cardId string, userID string) (*dao.KfUser, error) {
	key := d.getKey(cardId, userID)
	user := getCacheByKey[dao.KfUser](ctx, key)
	if user != nil {
		return user, nil
	}
	// 查数据库.
	var repo repository.KFUserRepository
	user, ok, err := repo.FindByToken(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, err
	}
	d.SetDBUser(ctx, cardId, user)
	return user, nil
}
