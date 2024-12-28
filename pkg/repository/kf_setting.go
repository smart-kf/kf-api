package repository

import (
	"context"
	"errors"
	xlogger "github.com/clearcodecn/log"
	"github.com/smart-fm/kf-api/pkg/db"
	"gorm.io/gorm"
)

type KFSettingRepository struct{}

func (r *KFSettingRepository) GetByCardID(ctx context.Context, cardID string) (*db.KFSettings, bool, error) {
	tx := db.GetDBFromContext(ctx)
	var setting db.KFSettings
	res := tx.Where("card_id = ?", cardID).First(&setting)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &setting, true, nil
}

// SaveOne save一条数据
func (r *KFCardRepository) SaveOne(ctx context.Context, setting *db.KFSettings) error {
	tx := db.GetDBFromContext(ctx)
	res := tx.Where("card_id = ?", setting.CardID).Save(setting)
	if res.Error != nil {
		xlogger.Error(ctx, "SaveOne-failed", xlogger.Err(res.Error))
		return res.Error
	}
	return nil
}
