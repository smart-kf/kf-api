package caches

import (
	"sync"

	"github.com/smart-fm/kf-api/infrastructure/redis"
)

var (
	BillSettingCacheInstance *BillSettingCache
	cacheOnce                = sync.Once{}
	CaptchaCacheInstance     *CaptchaCache
)

func InitCacheInstances() {
	cacheOnce.Do(
		func() {
			BillSettingCacheInstance = &BillSettingCache{}
			CaptchaCacheInstance = NewCaptchaCache(redis.GetRedisClient())
		},
	)
}
