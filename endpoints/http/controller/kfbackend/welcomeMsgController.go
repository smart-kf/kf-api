package kfbackend

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type WelcomeMsgController struct {
	BaseController
}

func (c *WelcomeMsgController) Upsert(ctx *gin.Context) {
	var req kfbackend.UpsertWelcomeMsgRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()

	db := mysql.GetDBFromContext(reqCtx)
	cardId := common.GetKFCardID(reqCtx)

	var model = dao.KfWelcomeMessage{
		Model: gorm.Model{
			ID: uint(req.Id),
		},
		CardId:  cardId,
		Content: req.Content,
		Type:    req.Type,
		Sort:    req.Sort,
		Enable:  req.Enable,
	}

	var err error
	if req.Id > 0 {
		err = db.Where("id = ?", req.Id).Save(&model).Error
	} else {
		var cnt int64
		err := db.Model(&dao.KfWelcomeMessage{}).Where("card_id = ?", cardId).Count(&cnt).Error
		if err != nil {
			c.Error(ctx, err)
			return
		}
		if cnt > 10 {
			c.Error(ctx, xerrors.NewCustomError("最多设置10条欢迎语"))
			return
		}
		err = db.Create(&model).Error
	}

	if err != nil {
		c.Error(ctx, err)
		return
	}

	c.Success(ctx, nil)
}

func (c *WelcomeMsgController) ListAll(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	db := mysql.GetDBFromContext(reqCtx)
	cardId := common.GetKFCardID(reqCtx)

	var data []*dao.KfWelcomeMessage
	db.Where("card_id = ?", cardId).Order("sort asc").Find(&data)

	var rsp []*kfbackend.KfWelcomeMessageResp
	for _, item := range data {
		rsp = append(
			rsp, &kfbackend.KfWelcomeMessageResp{
				Id:      int64(item.ID),
				Content: item.Content,
				Type:    item.Type,
				Sort:    item.Sort,
				Enable:  item.Enable,
			},
		)
	}

	c.Success(ctx, rsp)
}

func (c *WelcomeMsgController) Delete(ctx *gin.Context) {
	var req kfbackend.DeleteWelcomeRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()
	db := mysql.GetDBFromContext(reqCtx)

	cardId := common.GetKFCardID(reqCtx)
	n := db.Where("id = ? and card_id = ?", req.Id, cardId).Delete(&dao.KfWelcomeMessage{}).RowsAffected

	if n == 0 {
		c.Error(ctx, xerrors.NewCustomError("数据不存在"))
		return
	}

	c.Success(ctx, nil)
	return
}
