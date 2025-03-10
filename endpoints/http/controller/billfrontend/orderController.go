package billfrontend

import (
	"context"
	"fmt"
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/endpoints/http/vo/billfront"
	"github.com/smart-fm/kf-api/endpoints/http/vo/billfrontend"
	"github.com/smart-fm/kf-api/endpoints/nsq/producer"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	sdk "github.com/smart-fm/kf-api/pkg/usdtpayment"
	"github.com/smart-fm/kf-api/pkg/utils"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type OrderController struct {
	BaseController
}

// CreateOrder 卡密出售, 下单接口.
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req billfrontend.CreateOrderRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()
	// TODO:: 频繁下单限流校验.

	// 查询套餐.
	cardPackage, ok := config.GetConfig().GetPackageByID(req.PackageId)

	if !ok {
		c.Error(ctx, xerrors.NewCustomError("套餐不存在"))
		return
	}

	tx, reqCtx := mysql.Begin(reqCtx)
	defer tx.Rollback()

	var orderRepository repository.BillOrderRepository
	orderNo, err := caches.IdAtomicCacheInstance.GetBizId(reqCtx)
	if err != nil {
		c.Error(ctx, xerrors.NewCustomError("系统错误"))
		return
	}
	no := fmt.Sprintf("N%d", orderNo)

	var orderExpireDelay = config.GetConfig().BillConfig.OrderExpireTime
	var orderExpireTime = time.Now().UnixMicro() + orderExpireDelay*1000 // 使用毫秒.

	conf := config.GetConfig().Payment
	client := sdk.NewUsdtPaymentClient(conf.Host, conf.Token, 30*time.Second)

	// 1. 创建订单.
	var order = dao.Orders{
		Model: gorm.Model{
			ID: uint(orderNo),
		},
		PackageId:   cardPackage.Id,
		PackageDay:  cardPackage.Day,
		OrderNo:     no,
		Price:       cardPackage.Price * 1e6,     // 1 后面6个0
		Status:      constant.OrderStatusCreated, //
		ExpireTime:  orderExpireTime,
		Ip:          utils.ClientIP(ctx),
		Email:       req.Email,
		FromAddress: req.FromAddress,
	}

	createOrderResp, err := client.CreateOrder(
		reqCtx, &sdk.CreateOrderRequest{
			AppId:       conf.AppId,
			OrderId:     no,
			Name:        fmt.Sprintf("客服系统-%d日套餐", cardPackage.Day),
			Amount:      order.Price,
			FromAddress: order.FromAddress,
			Expire:      int(config.GetConfig().BillConfig.OrderExpireTime),
		},
	)

	if err != nil {
		c.Error(ctx, err)
		xlogger.Error(reqCtx, "orderCreateOne-usdtpayment-CreateOrder-失败:"+no, xlogger.Err(err))
		return
	}
	order.TradeId = createOrderResp.TradeId
	if err := orderRepository.CreateOne(reqCtx, &order); err != nil {
		c.Error(ctx, err)
		xlogger.Error(reqCtx, "orderCreateOne-失败:"+no, xlogger.Err(err))
		return
	}

	var (
		topic = config.GetConfig().NSQ.OrderExpireTopic
	)

	if err := producer.NSQProducer.DeferredPublish(
		topic, time.Duration(orderExpireDelay)*time.Second,
		[]byte(order.OrderNo),
	); err != nil {
		c.Error(ctx, err)
		xlogger.Error(reqCtx, "createOrder-DeferredPublish-失败:"+no, xlogger.Err(err))
		return
	}

	// 3. 返回订单号.
	if err := tx.Commit().Error; err != nil {
		c.Error(ctx, err)
		xlogger.Error(reqCtx, "createOrder-Commit-失败:"+no, xlogger.Err(err))
		return
	}

	c.Success(
		ctx, billfrontend.CreateOrderResponse{
			PaymentUrl: createOrderResp.PayUrl,
		},
	)
}

// Notify 订单 - 异步通知接口.
func (c *OrderController) Notify(ctx *gin.Context) {
	// TODO:: 目前是直接设置成 success, 并走购买成功逻辑
	var req billfront.OrderNotifyRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	if req.Status == 1 {
		c.Success(ctx, nil)
		return
	}
	reqCtx := ctx.Request.Context()

	tx, reqCtx := mysql.Begin(reqCtx)
	defer tx.Rollback()

	var orderRepository repository.BillOrderRepository
	order, ok, err := orderRepository.GetOrderByOrderNo(reqCtx, req.OrderId)
	if err != nil {
		xlogger.Error(reqCtx, "查询订单失败", xlogger.Err(err))
		c.Error(ctx, err)
		return
	}

	if !ok {
		xlogger.Error(reqCtx, "订单不存在")
		c.Error(ctx, xerrors.NewCustomError("订单不存在"))
		return
	}

	// 如果订单状态是已支付, 则不处理
	if order.Status == constant.OrderStatusPay {
		c.Success(ctx, nil)
		return
	}

	if req.Status == 3 {
		// 失败
		order.Status = constant.OrderStatusCancel
		order.ConfirmTime = time.Now().Unix()
		if err := tx.Where("id = ?", order.ID).Select(
			"status",
			"confirm_time",
		).Updates(order).Error; err != nil {
			xlogger.Error(reqCtx, "UpdateOrder failed", xlogger.Any("order", order), xlogger.Err(err))
			return
		}
		xlogger.Info(reqCtx, "收到订单失败回调", xlogger.Any("req", req))
		tx.Commit()
		c.Success(ctx, nil)
		return
	}

	order.TradeId = req.TradeId
	order.ConfirmTime = req.Timestamp
	order.Status = constant.OrderStatusPay
	order.ToAddress = req.Address

	// 查询一个卡密分配给他.
	var cardRepo repository.KFCardRepository
	card, exist, err := cardRepo.FindFirstCardByDay(reqCtx, order.PackageDay)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	if !exist {
		xlogger.Error(reqCtx, "告警-订单已支付但是库存不足", xlogger.Any("order", order))
		c.Error(ctx, xerrors.NewCustomError("库存不足"))
		return
	}

	// 更新card 数据.
	card.ExpireTime = time.Now().Unix() + int64(card.Day*86400) // 过期时间更新.
	card.SaleStatus = constant.SaleStatusSold

	// 更新订单数据
	order.CardID = card.CardID // 卡密id

	if err := tx.Where("id = ?", order.ID).Select(
		"card_id",
		"status",
		"confirm_time",
		"trade_id",
		"pay_usdt_address",
	).Updates(order).Error; err != nil {
		xlogger.Error(reqCtx, "UpdateOrder failed", xlogger.Any("order", order), xlogger.Err(err))
		c.Error(ctx, xerrors.NewCustomError("库存不足"))
		return
	}

	if err := tx.Where("id = ?", card.ID).Select(
		"sale_status",
		"expire_time",
	).Updates(card).Error; err != nil {
		xlogger.Error(reqCtx, "UpdateCard failed", xlogger.Any("order", order), xlogger.Err(err))
		c.Error(ctx, xerrors.NewCustomError("UpdateCard failed"))
		return
	}

	// 分配前台域名
	// 查找公共域名.
	var domainRepo repository.BillDomainRepository
	domain, exist, err := domainRepo.FindFirstPublic(reqCtx)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	if !exist {
		xlogger.Error(reqCtx, "告警-公共域名未找到", xlogger.Any("order", order))
		c.Error(ctx, xerrors.NewCustomError("公共域名未找到"))
		return
	}

	// 分配公共域名
	qrCodeDomain := dao.KFQRCodeDomain{
		Version:          1,
		ChangeQRCodeTime: time.Now().Unix(),
		Status:           constant.QRCodeNormal,
		Path:             "/s/" + utils.RandomPath(),
		CardID:           card.CardID,
		Domain:           domain.TopName,
		DomainId:         int64(domain.ID),
		IsPrivate:        false,
	}

	if err := tx.Create(&qrCodeDomain).Error; err != nil {
		xlogger.Error(
			reqCtx,
			"Create-KFQRCodeDomain failed",
			xlogger.Any("order", order),
			xlogger.Any("qrCodeDomain", qrCodeDomain),
			xlogger.Err(err),
		)
		c.Error(ctx, xerrors.NewCustomError("Create-KFQRCodeDomain failed"))
		return
	}

	// 设置域名绑定.
	if !domain.IsBind {
		domain.IsBind = true
		if err := tx.Where("id = ?", domain.ID).Select("is_bind").Updates(domain).Error; err != nil {
			xlogger.Error(
				reqCtx,
				"Create-UpdateDomainBind failed",
				xlogger.Any("order", order),
				xlogger.Any("domain", domain),
				xlogger.Err(err),
			)
			c.Error(ctx, xerrors.NewCustomError("Create-UpdateDomainBind failed"))
		}
	}

	// 分配默认系统配置.
	kfSetting := dao.NewDefaultKFSettings(card.CardID)
	if err := tx.Create(&kfSetting).Error; err != nil {
		xlogger.Error(
			reqCtx,
			"Create-KFSetting failed",
			xlogger.Any("order", order),
			xlogger.Err(err),
			xlogger.Any("setting", kfSetting),
		)
		c.Error(ctx, xerrors.NewCustomError("Create-KFSetting failed"))
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.Error(ctx, err)
		return
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {

			}
			ctx := context.Background()
			conf := config.GetConfig().Payment
			paymentClient := sdk.NewUsdtPaymentClient(conf.Host, conf.Token, 60*time.Second)
			paymentClient.SendMail(
				ctx, &sdk.SendMailRequest{
					From:       conf.FromMail,
					To:         order.Email,
					HtmlString: "支付成功: " + order.CardID,
				},
			)
		}()
	}()

	c.Success(ctx, order)
}
