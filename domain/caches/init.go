package caches

import (
	"sync"

	"github.com/smart-fm/kf-api/infrastructure/redis"
)

var (
	BillSettingCacheInstance *BillSettingCache
	cacheOnce                = sync.Once{}
	CaptchaCacheInstance     *CaptchaCache
	KfCardCacheInstance      *kfCardCacheInstance
	UserUnReadCacheInstance  *userUnReadCache
	KfUserExtraCacheInstance *kfUserExtraCache
	UserOnLineCacheInstance  *userOnLineCache
	ImSessionCacheInstance   *imSessionCache
	KfUserCacheInstance      *kfUserCache
	KfAuthCacheInstance      *kfAuthCache
)

func InitCacheInstances() {
	cacheOnce.Do(
		func() {
			BillSettingCacheInstance = &BillSettingCache{}
			CaptchaCacheInstance = NewCaptchaCache(redis.GetRedisClient())
			KfCardCacheInstance = &kfCardCacheInstance{}
			UserUnReadCacheInstance = &userUnReadCache{}
			KfUserExtraCacheInstance = &kfUserExtraCache{}
			UserOnLineCacheInstance = &userOnLineCache{}
			ImSessionCacheInstance = &imSessionCache{}
			KfUserCacheInstance = &kfUserCache{}
			KfAuthCacheInstance = &kfAuthCache{}
		},
	)
}
