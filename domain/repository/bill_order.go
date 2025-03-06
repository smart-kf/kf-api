package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type BillOrderRepository struct{}

// CreateOne 创建订单
func (r *BillOrderRepository) CreateOne(ctx context.Context, order *dao.Orders) error {
	tx := mysql.GetDBFromContext(ctx)
	if err := tx.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillOrderRepository) UpdateOrder(ctx context.Context, order *dao.Orders) error {
	tx := mysql.GetDBFromContext(ctx)
	version := order.Version
	order.Version++
	if err := tx.Where("id = ? and version = ?", order.ID, version).Updates(order).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillOrderRepository) GetOrderByID(ctx context.Context, id int64) (*dao.Orders, bool, error) {
	tx := mysql.GetDBFromContext(ctx)
	var order dao.Orders
	err := tx.Where("id = ?", id).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &order, true, nil
}

func (r *BillOrderRepository) GetOrderByOrderNo(ctx context.Context, orderNo string) (*dao.Orders, bool, error) {
	tx := mysql.GetDBFromContext(ctx)
	var order dao.Orders
	err := tx.Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &order, true, nil
}

type ListOrderOptions struct {
	common.PageRequest
	OrderNo     string `json:"orderNo" doc:"订单id"`
	TradeId     string `json:"tradeId" doc:"区块链上交易id"`
	FromAddress string `json:"fromAddress" doc:"客户地址"`
	ToAddress   string `json:"toAddress" doc:"接收地址"`
	Status      int    `json:"status" doc:"1=等待支付,2=支付成功,3=已取消"`
}

func (r *BillOrderRepository) ListOrder(ctx context.Context, req ListOrderOptions) ([]*dao.Orders, int64, error) {
	tx := mysql.GetDBFromContext(ctx)
	if req.FromAddress != "" {
		tx = tx.Where("from_address = ?", req.FromAddress)
	}
	if req.ToAddress != "" {
		tx = tx.Where("to_address = ?", req.ToAddress)
	}
	if req.TradeId != "" {
		tx = tx.Where("trade_id = ?", req.TradeId)
	}
	if req.OrderNo != "" {
		tx = tx.Where("order_no = ?", req.OrderNo)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	return common.Paginate[*dao.Orders](tx, &req.PageRequest)
}
