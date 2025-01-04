package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/common"
)

type BillDomainRepository struct {
}

func (r *BillDomainRepository) FindByTopName(ctx context.Context, topName string) (*dao.BillDomain, bool, error) {
	tx := mysql.DB()

	var domain dao.BillDomain
	err := tx.Where("top_name = ?", topName).First(&domain).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &domain, true, nil
}

func (r *BillDomainRepository) FindByID(ctx context.Context, id int64) (*dao.BillDomain, bool, error) {
	tx := mysql.DB()

	var domain dao.BillDomain
	err := tx.Where("id = ?", id).First(&domain).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &domain, true, nil
}

type ListDomainOption struct {
	*common.PageRequest
	IsPublic bool
	Status   int
	IsBind   bool
}

func (r *BillDomainRepository) List(ctx context.Context, options *ListDomainOption) ([]*dao.BillDomain, int64, error) {
	tx := mysql.GetDBFromContext(ctx)
	if options.IsPublic {
		tx = tx.Where("is_public = ?", options.IsPublic)
	}
	if options.IsBind {
		tx = tx.Where("is_bind = ?", options.IsBind)
	}
	if options.Status != 0 {
		tx = tx.Where("status = ?", options.Status)
	}
	return common.Paginate[*dao.BillDomain](tx, options.PageRequest)
}

func (r *BillDomainRepository) DeleteByID(ctx context.Context, id int64) error {
	tx := mysql.GetDBFromContext(ctx)
	return tx.Where("id = ?", id).Delete(&dao.BillDomain{}).Error
}
