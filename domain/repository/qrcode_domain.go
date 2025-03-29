package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type QRCodeDomainRepository struct{}

// FindByPath TODO:: cache
func (r *QRCodeDomainRepository) FindByPath(ctx context.Context, path string) (*dao.KFQRCodeDomain, bool, error) {
	tx := mysql.GetDBFromContext(ctx)
	var data dao.KFQRCodeDomain
	if err := tx.Where("path = ?", path).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &data, true, nil
}

// FindByCardID TODO:: cache
func (r *QRCodeDomainRepository) FindByCardID(ctx context.Context, cardID string) (*dao.KFQRCodeDomain, bool, error) {
	tx := mysql.GetDBFromContext(ctx)
	var data dao.KFQRCodeDomain
	if err := tx.Where("card_id = ?", cardID).Order("version desc").First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &data, true, nil
}

// FindByCardID
func (r *QRCodeDomainRepository) FindByIdAndCardID(ctx context.Context, id int64, cardID string) (
	*dao.KFQRCodeDomain, bool,
	error,
) {
	tx := mysql.GetDBFromContext(ctx)
	var data dao.KFQRCodeDomain
	if err := tx.Where("id = ? and card_id = ?", id, cardID).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &data, true, nil
}

// FindDomain
func (r *QRCodeDomainRepository) FindDomain(ctx context.Context, cardId string) ([]*dao.KFQRCodeDomain, error) {
	tx := mysql.GetDBFromContext(ctx)
	var data []*dao.KFQRCodeDomain
	sql := `SELECT *
FROM (SELECT *,
             ROW_NUMBER() OVER (PARTITION BY domain ORDER BY version DESC) AS rn
      FROM kf_qrcode_domain where card_id = ?) sub
WHERE rn = 1;`
	tx.Raw(sql, cardId).Find(&data)
	return data, nil
}

// CreateOne
func (r *QRCodeDomainRepository) CreateOne(ctx context.Context, qrCode *dao.KFQRCodeDomain) error {
	tx := mysql.GetDBFromContext(ctx)
	if err := tx.Create(qrCode).Error; err != nil {
		return err
	}
	return nil
}

func (r *QRCodeDomainRepository) DisableOld(ctx context.Context, cardId string, domain string, notId int64) error {
	tx := mysql.GetDBFromContext(ctx)
	if err := tx.Model(&dao.KFQRCodeDomain{}).Where(
		"card_id = ? and domain = ? and id != ?", cardId, domain,
		notId,
	).Update(
		"status",
		constant.QRCodeDisable,
	).Error; err != nil {
		return err
	}
	return nil
}

func (r *QRCodeDomainRepository) UpdateStatus(ctx context.Context, cardId string, id int64, status int) error {
	tx := mysql.GetDBFromContext(ctx)
	if err := tx.Model(&dao.KFQRCodeDomain{}).Where("card_id = ? and id = ?", cardId, id).Update(
		"status",
		status,
	).Error; err != nil {
		return err
	}
	return nil
}
