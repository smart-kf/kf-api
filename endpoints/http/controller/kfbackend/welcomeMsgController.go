package kfbackend

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
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
		MsgType: req.MsgType,
		Title:   req.Title,
	}

	var err error
	if req.Id > 0 {
		var exist dao.KfWelcomeMessage
		err := db.Where("id = ? and msg_type = ? and deleted_at is null", req.Id, req.MsgType).First(&exist).Error
		if err != nil {
			c.Error(ctx, err)
			return
		}
		exist.Content = req.Content
		exist.Type = req.Type
		exist.Enable = req.Enable
		exist.Sort = req.Sort
		exist.Title = req.Title

		err = db.Where("id = ? and msg_type = ?", req.Id, req.MsgType).Save(exist).Error
	} else {
		if req.MsgType == constant.WelcomeMsg {
			var cnt int64
			err := db.Model(&dao.KfWelcomeMessage{}).Where(
				"card_id = ? and msg_type = ? and deleted_at is null",
				cardId,
				req.MsgType,
			).Count(&cnt).Error
			if err != nil {
				c.Error(ctx, err)
				return
			}
			if cnt >= 10 {
				c.Error(ctx, xerrors.NewCustomError("最多设置10条欢迎语"))
				return
			}
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
	var req kfbackend.ListAllRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()
	db := mysql.GetDBFromContext(reqCtx)
	cardId := common.GetKFCardID(reqCtx)

	var data []*dao.KfWelcomeMessage
	var cnt int64
	tx := db.Where("card_id = ? and msg_type = ?", cardId, req.MsgType).Order("sort asc")
	tx = tx.Model(&dao.KfWelcomeMessage{}).Count(&cnt)
	if req.Page != nil && req.PageSize != nil {
		tx = tx.Limit(int(req.GetPage())).Offset(int((req.GetPage() - 1) * req.GetPageSize()))
	}
	tx.Find(&data)
	var rsp []*kfbackend.KfWelcomeMessageResp
	for _, item := range data {
		rsp = append(
			rsp, &kfbackend.KfWelcomeMessageResp{
				Id:      int64(item.ID),
				Content: item.Content,
				Type:    item.Type,
				Sort:    item.Sort,
				Enable:  item.Enable,
				Keyword: item.Keyword,
				Title:   item.Title,
			},
		)
	}

	c.Success(
		ctx, &kfbackend.KfWelcomeMessageListResp{
			List:  rsp,
			Total: cnt,
		},
	)
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

func (c *WelcomeMsgController) CopyCardMsg(ctx *gin.Context) {
	var req kfbackend.CopyCardMsgRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()
	cardId := common.GetKFCardID(reqCtx)

	// 1. 查询对方的卡密.
	var cardRepo repository.KFCardRepository
	_, ok, err := cardRepo.FindByCardID(reqCtx, cardId)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	if !ok {
		c.Error(ctx, xerrors.NewParamsErrors("卡密不存在"))
		return
	}

	tx, txCtx := mysql.Begin(reqCtx)
	defer func() {
		tx.Rollback()
	}()

	var msgRepo repository.KfWelcomeMessage

	if req.QuickReply {
		err = msgRepo.CopyFromCard(
			txCtx, repository.CopyParams{
				FromCardId:           req.CardID,
				ToCardId:             cardId,
				ReplaceTargetContent: req.ReplaceTargetContent,
				ReplaceContent:       req.ReplaceContent,
				MsgType:              constant.QuickReply,
			},
		)
		if err != nil {
			c.Error(ctx, err)
			return
		}
	}
	if req.WelcomeMsg {
		err = msgRepo.CopyFromCard(
			txCtx, repository.CopyParams{
				FromCardId:           req.CardID,
				ToCardId:             cardId,
				ReplaceTargetContent: req.ReplaceTargetContent,
				ReplaceContent:       req.ReplaceContent,
				MsgType:              constant.WelcomeMsg,
			},
		)
		if err != nil {
			c.Error(ctx, err)
			return
		}
	}
	if req.SmartReply {
		err = msgRepo.CopyFromCard(
			txCtx, repository.CopyParams{
				FromCardId:           req.CardID,
				ToCardId:             cardId,
				ReplaceTargetContent: req.ReplaceTargetContent,
				ReplaceContent:       req.ReplaceContent,
				MsgType:              constant.SmartReply,
			},
		)
		if err != nil {
			c.Error(ctx, err)
			return
		}
	}

	if req.Nickname || req.Avatar || req.Settings {
		var settingRepo repository.KFSettingRepository
		err = settingRepo.CopyFromCard(
			txCtx, repository.CopySettingParam{
				FromCardId: req.CardID,
				ToCardId:   cardId,
				Nickname:   req.Nickname,
				Avatar:     req.Avatar,
				Settings:   req.Settings,
			},
		)
		if err != nil {
			c.Error(ctx, err)
			return
		}
		caches.KfSettingCache.DeleteOne(ctx, cardId)
	}

	tx.Commit()
	c.Success(ctx, nil)
}
