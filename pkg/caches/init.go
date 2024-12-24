package caches

import (
	"github.com/smart-fm/kf-api/pkg/db"
	"sync"
)

var (
	BillSettingCacheInstance *BillSettingCache
	cacheOnce                = sync.Once{}
	CaptchaCacheInstance     *CaptchaCache
)

func InitCacheInstances() {
	cacheOnce.Do(func() {
		BillSettingCacheInstance = &BillSettingCache{}
		CaptchaCacheInstance = NewCaptchaCache(db.GetRedisClient())
	})
}
