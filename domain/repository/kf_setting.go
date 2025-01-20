package repository

import (
	"context"

	xlogger "github.com/clearcodecn/log"

	mysql2 "github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type KFSettingRepository struct{}

func (r *KFSettingRepository) MustGetByCardID(ctx context.Context, cardID string) (*dao.KFSettings, error) {
	tx := mysql2.GetDBFromContext(ctx)
	var setting dao.KFSettings
	res := tx.Where("card_id = ?", cardID).First(&setting)
	if err := res.Error; err != nil {
		return dao.NewDefaultKFSettings(cardID), nil
	}
	return &setting, nil
}

// SaveOne save一条数据
func (r *KFSettingRepository) SaveOne(ctx context.Context, setting *dao.KFSettings) error {
	tx := mysql2.GetDBFromContext(ctx)
	res := tx.Where("card_id = ?", setting.CardID).Save(setting)
	if res.Error != nil {
		xlogger.Error(ctx, "SaveOne-failed", xlogger.Err(res.Error))
		return res.Error
	}
	return nil
}
