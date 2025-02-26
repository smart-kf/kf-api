package repository

import (
	"context"

	xlogger "github.com/clearcodecn/log"

	mysql2 "github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/xerrors"
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

type CopySettingParam struct {
	FromCardId string
	ToCardId   string
	Nickname   bool
	Avatar     bool
	Settings   bool
}

func (r *KFSettingRepository) CopyFromCard(ctx context.Context, param CopySettingParam) error {
	db := mysql2.GetDBFromContext(ctx)
	var setting dao.KFSettings
	err := db.Where("card_id = ?", param.FromCardId).First(&setting).Error
	if err != nil {
		return xerrors.NewParamsErrors("配置不存在")
	}
	var newSetting dao.KFSettings
	db.Where("card_id = ?", param.ToCardId).First(&newSetting)
	if param.Nickname {
		newSetting.Nickname = setting.Nickname
	}
	if param.Avatar {
		newSetting.AvatarURL = setting.AvatarURL
	}
	if param.Settings {
		newSetting.AppleFilter = setting.AppleFilter
		newSetting.WSFilter = setting.WSFilter
		newSetting.WechatFilter = setting.WechatFilter
		newSetting.IPProxyFilter = setting.IPProxyFilter
		newSetting.DeviceFilter = setting.DeviceFilter
		newSetting.SimulatorFilter = setting.SimulatorFilter
		newSetting.NewMessageVoice = setting.NewMessageVoice
		newSetting.QRCodeEnabled = setting.QRCodeEnabled
	}

	return db.Where("card_id = ?", newSetting.CardID).Save(&newSetting).Error
}
