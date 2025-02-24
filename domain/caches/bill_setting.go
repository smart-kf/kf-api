package caches

import (
	"context"
	"encoding/json"

	redis2 "github.com/redis/go-redis/v9"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/infrastructure/redis"
)

type BillSettingCache struct {
}

func (c *BillSettingCache) GetNotice() dao.Notice {
	ctx := context.Background()
	setting := c.getSetting(ctx)
	if setting == nil {
		return dao.Notice{}
	}
	return setting.Notice
}

func (c BillSettingCache) GetTestingCardMinute() int {
	ctx := context.Background()
	setting := c.getSetting(ctx)
	if setting == nil {
		return 15
	}
	return setting.TestingCardMinute
}

func (s *BillSettingCache) OneDayCardPrice() float64 {
	ctx := context.Background()
	setting := s.getSetting(ctx)
	if setting == nil {
		return 0
	}
	return setting.DailyPackage.Price
}

func (c *BillSettingCache) getSetting(ctx context.Context) *dao.BillSettingModel {
	cli := redis.GetRedisClient()
	data, err := cli.Get(ctx, "bill_setting").Bytes()
	if err != nil {
		if err == redis2.Nil {
			var repo repository.BillSettingRepository
			setting, err := repo.GetSetting(ctx)
			if err != nil {
				return nil
			}
			settingBytes, err := json.Marshal(setting)
			if err != nil {
				return nil
			}
			cli.SetEx(ctx, "bill_setting", settingBytes, 60*60*24)
			data = settingBytes
		}
	}

	var setting dao.BillSettingModel
	err = json.Unmarshal(data, &setting)
	if err != nil {
		return nil
	}
	return &setting
}

func (BillSettingCache) DeleteCache(ctx context.Context) {
	cli := redis.GetRedisClient()
	cli.Del(ctx, "bill_setting")
}
