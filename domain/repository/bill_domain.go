package repository

import (
	"context"
	"math/rand"

	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type BillDomainRepository struct {
}

func (r *BillDomainRepository) CountByTopNames(ctx context.Context, topName []string) (int64, error) {
	tx := mysql.DB()
	var cnt int64
	err := tx.Model(&dao.BillDomain{}).Where("top_name in ?", topName).Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
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

func (r *BillDomainRepository) FindFirstPublic(ctx context.Context) (*dao.BillDomain, bool, error) {
	tx := mysql.GetDBFromContext(ctx)

	var domain []*dao.BillDomain
	err := tx.Where("is_public = ? and status = ?", true, constant.DomainStatusNormal).Find(&domain).Error

	if err != nil {
		return nil, false, err
	}

	if len(domain) == 0 {
		return nil, false, nil
	}

	return domain[rand.Intn(len(domain))], true, nil
}

func (r *BillDomainRepository) CountPrivateDomain(ctx context.Context) (int64, error) {
	var cnt int64
	tx := mysql.GetDBFromContext(ctx)
	err := tx.Model(&dao.BillDomain{}).Where(
		"is_public = ? and status = ?", false,
		constant.DomainStatusNormal,
	).Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
