package bill

import (
	"github.com/gin-gonic/gin"
	"std-api/pkg/common"
	"std-api/pkg/constant"
	"std-api/pkg/db"
	"std-api/pkg/repository"
	"std-api/pkg/utils"
	"std-api/pkg/xerrors"
)

type BatchAddCardRequest struct {
	CardType constant.CardType `json:"cardType" binding:"required" validate:"required,oneof=1 2" doc:"卡密类型: 1=正式卡，2=测试卡"`
	Days     int               `json:"days" doc:"天数,正式卡必填,测试卡忽略"`
	Num      int               `json:"num" doc:"数量,1-100之间的整数" binding:"required" validate:"required,gte=1,lte=100"`
}

type BatchAddResponse struct {
	Num int `json:"num"`
}

func (req *BatchAddCardRequest) Validate() error {
	if req.CardType == constant.CardTypeNormal {
		if req.Days <= 0 {
			return xerrors.NewParamsErrors("请填写天数")
		}
	}
	return nil
}

type CardController struct {
	BaseController
}

func (c *CardController) BatchAddCard(ctx *gin.Context) {
	var req BatchAddCardRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	var cards []*db.KFCard
	for i := 0; i < req.Num; i++ {
		cards = append(cards, &db.KFCard{
			CardID:      utils.RandomCard(),
			SaleStatus:  constant.SaleStatusOnSale,
			LoginStatus: constant.LoginStatusUnLogin,
			CardType:    req.CardType,
			Day:         req.Days,
		})
	}

	var cardRepo repository.KFCardRepository
	if err := cardRepo.CreateBatch(ctx, cards); err != nil {
		c.Error(ctx, err)
		return
	}

	db.AddBillLog(common.GetBillAccount(ctx).Username, "BatchAddCard", "批量添加卡密")

	c.Success(ctx, len(cards))
}

type ListCardRequest struct {
	common.PageRequest
	SaleStatus         constant.SaleStatus  `json:"saleStatus" doc:"卡片状态,1=销售中，2=下架，3=已出售"`
	LoginStatus        constant.LoginStatus `json:"loginStatus" doc:"登录状态: 1=未登录过，2=登录过"`
	CardType           constant.CardType    `json:"cardType" doc:"卡片类型: 1正式卡, 2测试卡"`
	ExpireStartTime    int64                `json:"expireStartTime" doc:"过期时间-开始时间，秒"`
	ExpireEndTime      int64                `json:"expireEndTime" doc:"过期时间-结束时间，秒"`
	CardID             string               `json:"cardID" doc:"卡密id"`
	LastLoginTimeStart int64                `json:"lastLoginTimeStart" doc:"上次登录时间-开始，秒"`
	LastLoginTimeEnd   int64                `json:"lastLoginTimeEnd" doc:"上次登录时间-开始，秒"`
}

type ListCardResponse struct {
	List  []*KFCardResponse `json:"list" doc:"列表数据"`
	Total int64             `json:"total" doc:"统计"`
}

type KFCardResponse struct {
	ID            uint                 `json:"id" doc:"主键id"`
	CardID        string               `json:"cardId" doc:"卡密id"`
	Password      string               `json:"password" doc:"密码"`
	SaleStatus    constant.SaleStatus  `json:"saleStatus" doc:"销售状态"`
	LoginStatus   constant.LoginStatus `json:"loginStatus" doc:"登录状态"`
	CardType      constant.CardType    `json:"cardType" doc:"卡片类型"`
	Day           int                  `json:"day"  doc:"卡密的天数"`
	ExpireTime    int64                `json:"expireTime" doc:"过期时间"`
	LastLoginTime int64                `json:"lastLoginTime" doc:"上次登录时间"`
}

func (c *CardController) List(ctx *gin.Context) {
	var req ListCardRequest
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

	var rsp []*KFCardResponse
	for _, item := range res {
		rsp = append(rsp, &KFCardResponse{
			ID:            item.ID,
			CardID:        item.CardID,
			Password:      item.Password,
			SaleStatus:    item.SaleStatus,
			LoginStatus:   item.LoginStatus,
			CardType:      item.CardType,
			Day:           item.Day,
			ExpireTime:    item.ExpireTime,
			LastLoginTime: item.LastLoginTime,
		})
	}

	httpResp := ListCardResponse{
		List:  rsp,
		Total: cnt,
	}

	c.Success(ctx, &httpResp)
}

type UpdateStatusRequest struct {
	ID     uint                `json:"id" binding:"required" doc:"主键id" validate:"required,gt=0"`
	Status constant.SaleStatus `json:"status" binding:"required" doc:"卡片状态,1=销售中，2=下架，3=已出售" validate:"required,oneof=1 2 3"`
}

// UpdateStatus 修改卡片状态
// 1. 可以将卡片置为已出售，出售状态可以登录
func (c *CardController) UpdateStatus(ctx *gin.Context) {
	var req UpdateStatusRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	tx, txCtx := db.Begin(ctx)
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
	c.Success(ctx, nil)
	return
}
