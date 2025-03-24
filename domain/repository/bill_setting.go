package repository

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

const settingKey = "system_setting"

type BillSettingRepository struct {
}

func (r *BillSettingRepository) GetSetting(ctx context.Context) (*dao.BillSettingModel, error) {
	var setting dao.BillSetting
	db := mysql.GetDBFromContext(ctx)
	if err := db.Where("`key` = ?", settingKey).Model(&dao.BillSetting{}).First(&setting).Error; err != nil {
		return nil, err
	}
	var model dao.BillSettingModel
	err := json.Unmarshal([]byte(setting.Value), &model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *BillSettingRepository) UpsertSettings(
	ctx context.Context, setting *dao.BillSettingModel,
	updateIfExist bool,
) error {
	settingValue, err := json.Marshal(setting)
	if err != nil {
		return err
	}
	db := mysql.GetDBFromContext(ctx)

	var settingModel dao.BillSetting
	err = db.Where("`key` = ?", settingKey).First(&settingModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(
				&dao.BillSetting{
					Key:   settingKey,
					Value: string(settingValue),
				},
			).Error
		}
	} else {
		if updateIfExist {
			err = db.Model(&dao.BillSetting{}).Where("`key` = ?", settingKey).
				Update("value", string(settingValue)).Error
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *BillSettingRepository) InitDefault() {
	model := &dao.BillSettingModel{
		DailyPackage: dao.Package{
			Id:    "daily",
			Days:  1,
			Price: 15,
		},
		WeeklyPackage: dao.Package{
			Id:    "weekly",
			Days:  7,
			Price: 15 * 7 * 0.9,
		},
		MonthlyPackage: dao.Package{
			Id:    "monthly",
			Days:  7,
			Price: 15 * 30 * 0.8,
		},
		TestingCardMinute: 15,
		Payment: dao.Payment{
			PayUrl: "",
			Token:  "",
			AppId:  "kf",
			Email:  "",
		},
		Notice: dao.Notice{
			Content: "testing content",
			Enable:  false,
		},
		DomainPrice: 10,
	}

	r.UpsertSettings(context.Background(), model, false)
}
