package billfrontend

import (
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/endpoints/http/vo/billfrontend"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/infrastructure/redis"
	"github.com/smart-fm/kf-api/pkg/utils"
)

type OrderController struct {
	BaseController
}

// CreateOrder 卡密出售, 下单接口.
func (c *BaseController) CreateOrder(ctx *gin.Context) {
	var req billfrontend.CreateOrderRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()
	// TODO:: 频繁下单限流校验.

	no := utils.RandomOrderNo()
	tx, reqCtx := mysql.Begin(reqCtx)
	defer tx.Rollback()

	var orderRepository repository.BillOrderRepository
	for {
		// 检测订单号是否存在.
		_, ok, err := orderRepository.GetOrderByOrderNo(reqCtx, no)
		if err != nil {
			xlogger.Error(reqCtx, "GetOrderByOrderNo-失败:"+no, xlogger.Err(err))
			c.Error(ctx, err)
			return
		}
		if !ok {
			break
		}
	}

	// 1. 创建订单.
	var order = dao.Orders{
		CardID:         "",
		OrderNo:        no,
		PayUsdtAddress: "",                          // TODO:: 支付地址.
		Price:          15 * 1e6,                    // TODO:: payment
		Status:         constant.OrderStatusCreated, //
		ConfirmTime:    0,
		Ip:             utils.ClientIP(ctx),
		Area:           "", // TODO:: ip2region
		Version:        1,
		ExpireTime:     time.Now().UnixMicro() + config.GetConfig().BillConfig.OrderExpireTime*1000, // 使用毫秒.
	}

	if err := orderRepository.CreateOne(reqCtx, &order); err != nil {
		c.Error(ctx, err)
		xlogger.Error(reqCtx, "orderCreateOne-失败:"+no, xlogger.Err(err))
		return
	}

	// 2. 设置取消队列.
	redisClient := redis.GetRedisClient()
	if err := utils.ZAdd(
		reqCtx,
		redisClient,
		constant.OrderExpireZSetKey,
		utils.ZSetMember{Member: no, Score: order.ExpireTime},
	); err != nil {
		c.Error(ctx, err)
		xlogger.Error(reqCtx, "createOrder-ZAdd-失败:"+no, xlogger.Err(err))
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
			OrderNo: order.OrderNo,
		},
	)
}
