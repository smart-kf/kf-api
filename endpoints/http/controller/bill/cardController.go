package bill

import (
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/endpoints/cron/billlog"
	"github.com/smart-fm/kf-api/endpoints/http/middleware"
	"github.com/smart-fm/kf-api/endpoints/http/vo/bill"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/utils"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type CardController struct {
	BaseController
}

func (c *CardController) BatchAddCard(ctx *gin.Context) {
	var req bill.BatchAddCardRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	var cards []*dao.KFCard
	oneDayCardPrice := caches.BillSettingCacheInstance.OneDayCardPrice()
	price := int64(oneDayCardPrice * req.Days)
	for i := 0; i < req.Num; i++ {
		cards = append(
			cards, &dao.KFCard{
				CardID:      utils.RandomCard(10),
				SaleStatus:  constant.SaleStatusOnSale,
				LoginStatus: constant.LoginStatusUnLogin,
				CardType:    req.CardType,
				Day:         req.Days,
				Price:       price,
			},
		)
	}

	var cardRepo repository.KFCardRepository
	if err := cardRepo.CreateBatch(ctx, cards); err != nil {
		c.Error(ctx, err)
		return
	}

	billlog.AddBillLog(middleware.GetBillAccount(ctx).Username, "BatchAddCard", "批量添加卡密")

	c.Success(ctx, len(cards))
}

func (c *CardController) List(ctx *gin.Context) {
	var req bill.ListCardRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	var listOptions = repository.ListCardOption{
		CardID:             req.CardID,
		CardType:           req.CardType,
		SaleStatus:         req.SaleStatus,
		LoginStatus:        req.LoginStatus,
		ExpireStart:        req.ExpireStartTime,
		ExpireEnd:          req.ExpireEndTime,
		PageRequest:        &req.PageRequest,
		LastLoginTimeStart: req.LastLoginTimeStart,
		LastLoginTimeEnd:   req.LastLoginTimeEnd,
	}
	var cardRepo repository.KFCardRepository

	res, cnt, err := cardRepo.List(ctx.Request.Context(), &listOptions)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	var rsp []*bill.KFCardResponse
	for _, item := range res {
		rsp = append(
			rsp, &bill.KFCardResponse{
				ID:            item.ID,
				CardID:        item.CardID,
				Password:      item.Password,
				SaleStatus:    item.SaleStatus,
				LoginStatus:   item.LoginStatus,
				CardType:      item.CardType,
				Day:           item.Day,
				ExpireTime:    item.ExpireTime,
				LastLoginTime: item.LastLoginTime,
			},
		)
	}

	httpResp := bill.ListCardResponse{
		List:  rsp,
		Total: cnt,
	}

	c.Success(ctx, &httpResp)
}

// UpdateStatus 修改卡片状态
// 1. 可以将卡片置为已出售，出售状态可以登录
func (c *CardController) UpdateStatus(ctx *gin.Context) {
	var req bill.UpdateStatusRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	tx, txCtx := mysql.Begin(ctx)
	defer tx.Rollback()

	var cardRepo repository.KFCardRepository
	card, ok, err := cardRepo.GetByID(txCtx, req.ID)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	if !ok {
		c.Error(ctx, xerrors.NewParamsErrors("数据不存在"))
		return
	}
	if card.SaleStatus == req.Status {
		c.Success(ctx, nil)
		return
	}
	card.SaleStatus = req.Status

	err = cardRepo.UpdateOne(txCtx, card)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	tx.Commit()
	c.Success(ctx, nil)
	return
}
