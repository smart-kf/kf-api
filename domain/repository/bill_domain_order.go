package repository

import (
	"context"

	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type BillDomainOrderRepository struct{}

// CreateOne 创建订单
func (r *BillDomainOrderRepository) CreateOne(ctx context.Context, order *dao.DomainOrder) error {
	tx := mysql.GetDBFromContext(ctx)
	if err := tx.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillDomainOrderRepository) FindByOrderNo(ctx context.Context, orderNo string) (*dao.DomainOrder, error) {
	tx := mysql.GetDBFromContext(ctx)
	var order dao.DomainOrder
	if err := tx.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *BillDomainOrderRepository) UpdateStatus(ctx context.Context, order *dao.DomainOrder) error {
	tx := mysql.GetDBFromContext(ctx)
	if err := tx.Model(order).Where("id = ?", order.ID).Update("status", order.Status).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillDomainOrderRepository) CheckUnPayedByCardId(ctx context.Context, cardId string) (bool, error) {
	tx := mysql.GetDBFromContext(ctx)
	var cnt int64
	if err := tx.Model(&dao.DomainOrder{}).Where(
		"card_id = ? and status = ?", cardId,
		constant.OrderStatusCreated,
	).Count(&cnt).Error; err != nil {
		return false, err
	}
	if cnt > 0 {
		return true, nil
	}
	return false, nil
}
