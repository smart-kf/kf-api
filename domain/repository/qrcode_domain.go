package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type QRCodeDomainRepository struct{}

// FindByPath TODO:: cache
func (r *QRCodeDomainRepository) FindByPath(ctx context.Context, path string) (*dao.KFQRCode, bool, error) {
	tx := mysql.GetDBFromContext(ctx)
	var data dao.KFQRCode
	if err := tx.Where("path = ?", path).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &data, true, nil
}

// FindByCardID TODO:: cache
func (r *QRCodeDomainRepository) FindByCardID(ctx context.Context, cardID string) (*dao.KFQRCode, bool, error) {
	tx := mysql.GetDBFromContext(ctx)
	var data dao.KFQRCode
	if err := tx.Where("card_id = ?", cardID).Order("version desc").First(&data).Error; err != nil {
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
	if err := tx.Where("card_id = ?", cardId).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// CreateOne
func (r *QRCodeDomainRepository) CreateOne(ctx context.Context, qrCode *dao.KFQRCode) error {
	tx := mysql.GetDBFromContext(ctx)
	if err := tx.Create(qrCode).Error; err != nil {
		return err
	}
	return nil
}
