package orders

import (
	"context"

	xlogger "github.com/clearcodecn/log"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
)

func ExpireDomainOrder(ctx context.Context, orderNo string) {
	var domainRepo repository.BillDomainOrderRepository
	order, err := domainRepo.FindByOrderNo(ctx, orderNo)
	if err != nil {
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
	if err := domainRepo.UpdateStatus(newCtx, order); err != nil {
		xlogger.Error(ctx, "更新订单状态失败", xlogger.Any("orderNo", orderNo), xlogger.Err(err))
		return
	}

	// 释放域名.
	var billDomainRepo repository.BillDomainRepository
	domain, ok, err := billDomainRepo.FindByID(newCtx, int64(order.DomainId))
	if err != nil {
		xlogger.Error(ctx, "更新订单状态失败", xlogger.Any("orderNo", orderNo), xlogger.Err(err))
		return
	}
	if !ok {
		xlogger.Error(ctx, "域名不存在", xlogger.Any("orderNo", orderNo))
		return
	}
	domain.Status = constant.DomainStatusNormal

	if err := billDomainRepo.Update(newCtx, domain); err != nil {
		xlogger.Error(ctx, "释放域名失败", xlogger.Any("orderNo", orderNo), xlogger.Err(err))
		return
	}

	tx.Commit()

	xlogger.Info(ctx, "域名订单过期处理成功", xlogger.Any("orderNo", orderNo))
}
