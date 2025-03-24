package orders

import (
	"context"

	xlogger "github.com/clearcodecn/log"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
)

func CardOrderExpire(ctx context.Context, orderNo string) {
	var orderRepo repository.BillOrderRepository
	order, ok, err := orderRepo.GetOrderByOrderNo(ctx, orderNo)
	if err != nil {
		return
	}
	if !ok {
		xlogger.Error(ctx, "订单不存在", xlogger.Any("orderNo", orderNo))
		return
	}
	// 查询订单状态.
	if order.Status != constant.OrderStatusCreated {
		xlogger.Info(ctx, "订单状态不是已创建状态，不需要处理", xlogger.Any("orderNo", orderNo))
		return
	}

	tx, newCtx := mysql.Begin(ctx)

	defer func() {
		tx.Rollback()
	}()

	// 更新订单状态.
	order.Status = constant.OrderStatusCancel
	if err := orderRepo.UpdateOrder(newCtx, order); err != nil {
		xlogger.Error(ctx, "更新订单状态失败", xlogger.Any("orderNo", orderNo), xlogger.Err(err))
		return
	}
	tx.Commit()
	xlogger.Info(ctx, "卡密订单过期处理成功", xlogger.Any("orderNo", orderNo))
}
