package bill

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"std-api/pkg/common"
	"std-api/pkg/constant"
	"std-api/pkg/db"
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
			Password:    "",
			SaleStatus:  constant.SaleStatusOnSale,
			LoginStatus: constant.LoginStatusUnLogin,
			CardType:    req.CardType,
			Day:         req.Days,
		})
	}

	var tx = db.DB()
	if err := tx.Model(&db.KFCard{}).CreateInBatches(cards, len(cards)).Error; err != nil {
		xlogger.Error("BatchAddCard-failed", xlogger.Err(err))
		c.Error(ctx, err)
		return
	}

	db.AddBillLog(common.GetBillAccount(ctx).Username, "BatchAddCard", "批量添加卡密")

	c.Success(ctx, len(cards))
}
