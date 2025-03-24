package orders

import (
	"context"
	"fmt"
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/samber/lo"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/endpoints/nsq/producer"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	sdk "github.com/smart-fm/kf-api/pkg/usdtpayment"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type CreateDomainDTO struct {
	FromAddress string
}

type CreateDomainResult struct {
	PayUrl  string `json:"payUrl"`
	TradeId string `json:"tradeId"`
}

func CreateDomainOrder(ctx context.Context, dto CreateDomainDTO) (*CreateDomainResult, error) {
	var domainRepo repository.BillDomainRepository
	cnt, err := domainRepo.CountPrivateDomain(ctx)
	if err != nil {
		return nil, err
	}

	// 检查有没有未支付的订单.
	var orderRepo repository.BillDomainOrderRepository
	ok, err := orderRepo.CheckUnPayedByCardId(ctx, common.GetKFCardID(ctx))
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, xerrors.NewCustomError("您存在未过期的订单，请处理后再试")
	}

	cardId := common.GetKFCardID(ctx)

	if cnt == 0 {
		return nil, xerrors.NewCustomError("没有可用的域名,请联系客服处理")
	}

	// 查询缓存.
	price := caches.BillSettingCacheInstance.GetDomainPrice()

	if price == 0 {
		return nil, xerrors.NewCustomError("系统配置错误，请联系客服处理")
	}

	tx, newCtx := mysql.Begin(ctx)

	defer func() {
		tx.Rollback()
	}()

	_ = cardId
	_, _ = tx, newCtx

	// 生成订单号.
	orderNo, err := caches.IdAtomicCacheInstance.GetBizId(newCtx)
	if err != nil {
		return nil, xerrors.NewCustomError("系统错误")
	}
	no := fmt.Sprintf("X%d", orderNo)
	_ = no

	// 1. 创建订单
	var order = &dao.DomainOrder{
		CardID:      cardId,
		OrderNo:     no,
		ToAddress:   "",
		FromAddress: dto.FromAddress,
		Price:       int64(price) * 1e6,
		Status:      constant.OrderStatusCreated,
		ConfirmTime: 0,
		ExpireTime:  time.Now().Add(constant.DomainExpireTime).Unix(), // 5分钟过期
	}
	// 2. 锁定域名
	domain, ok, err := domainRepo.LockOnePrivateDomain(newCtx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, xerrors.NewCustomError("没有可用的域名,请联系客服处理")
	}
	domain.Status = constant.DomainStatusLocked
	err = domainRepo.Update(newCtx, domain)
	if err != nil {
		return nil, err
	}
	// 4. 创建订单.
	order.Domain = domain.TopName
	order.DomainId = int(domain.ID)
	var (
		topic = config.GetConfig().NSQ.OrderExpireTopic
	)
	// 3. 推送取消订单事件.
	if err := producer.NSQProducer.DeferredPublish(
		topic, constant.DomainExpireTime,
		[]byte(order.OrderNo),
	); err != nil {
		xlogger.Error(ctx, "createOrder-DeferredPublish-失败:"+no, xlogger.Err(err))
		return nil, err
	}
	// 4. 创建订单.
	conf := config.GetConfig().Payment
	client := sdk.NewUsdtPaymentClient(conf.Host, conf.Token, 30*time.Second)
	createOrderResp, err := client.CreateOrder(
		newCtx, &sdk.CreateOrderRequest{
			AppId:       conf.AppId,
			OrderId:     no,
			Name:        "域名购买",
			Amount:      order.Price,
			FromAddress: order.FromAddress,
			Expire:      int(constant.DomainExpireTime / time.Second),
		},
	)
	if err != nil {
		xlogger.Error(newCtx, "orderCreateOne-usdtpayment-CreateOrder-失败:"+no, xlogger.Err(err))
		return nil, err
	}
	order.TradeId = createOrderResp.TradeId
	order.PayUrl = createOrderResp.PayUrl

	if err := orderRepo.CreateOne(newCtx, order); err != nil {
		xlogger.Error(newCtx, "orderCreateOne-usdtpayment-CreateOne-失败:"+no, xlogger.Err(err))
		return nil, err
	}
	tx.Commit()
	// 5. 返回订单号

	xlogger.Info(newCtx, "创建域名订单成功", xlogger.Any("order", order))
	// 6. 支付.
	return &CreateDomainResult{
		PayUrl:  createOrderResp.PayUrl,
		TradeId: createOrderResp.TradeId,
	}, nil
}

func ListOrder(ctx context.Context, cardId string) (*kfbackend.DomainOrderList, error) {
	var orderRepo repository.BillDomainOrderRepository
	_, _ = orderRepo, cardId
	// 1. 查询订单.
	// 2. 返回订单列表.
	var domain []*dao.DomainOrder
	db := mysql.GetDBFromContext(ctx)
	err := db.Where("card_id = ?", cardId).Order("id desc").Limit(30).Offset(0).Find(&domain).Error
	if err != nil {
		return nil, err
	}

	var (
		resp []*kfbackend.DomainOrder
	)
	lo.ForEach(
		domain, func(item *dao.DomainOrder, index int) {
			domain := item.Domain
			if item.Status != constant.OrderStatusPay {
				domain = ""
			}
			resp = append(
				resp, &kfbackend.DomainOrder{
					OrderNo:     item.OrderNo,
					ToAddress:   item.ToAddress,
					FromAddress: item.FromAddress,
					Price:       item.Price,
					Status:      item.Status,
					ConfirmTime: item.ConfirmTime,
					ExpireTime:  item.ExpireTime,
					Domain:      domain,
					TradeId:     item.TradeId,
					PayUrl:      item.PayUrl,
				},
			)
		},
	)

	return &kfbackend.DomainOrderList{
		List: resp,
	}, nil
}
