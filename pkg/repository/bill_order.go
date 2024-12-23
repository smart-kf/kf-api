package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"std-api/pkg/db"
)

type BillOrderRepository struct{}

// CreateOne 创建订单
func (r *BillOrderRepository) CreateOne(ctx context.Context, order *db.Orders) error {
	tx := db.GetDBFromContext(ctx)
	if err := tx.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillOrderRepository) UpdateOrder(ctx context.Context, order *db.Orders) error {
	tx := db.GetDBFromContext(ctx)
	version := order.Version
	order.Version++
	if err := tx.Where("id = ? and version = ?", order.ID, version).Updates(order).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillOrderRepository) GetOrderByID(ctx context.Context, id int64) (*db.Orders, bool, error) {
	tx := db.GetDBFromContext(ctx)
	var order db.Orders
	err := tx.Where("id = ?", id).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &order, true, nil
}

func (r *BillOrderRepository) GetOrderByOrderNo(ctx context.Context, orderNo string) (*db.Orders, bool, error) {
	tx := db.GetDBFromContext(ctx)
	var order db.Orders
	err := tx.Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &order, true, nil
}
