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
