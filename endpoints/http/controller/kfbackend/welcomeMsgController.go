package kfbackend

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/endpoints/cron/kflog"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/utils"
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
	repo := repository.KfWelcomeMessageRepository{}
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
		Keyword: req.Keyword,
	}
	err := repo.UpsertOne(reqCtx, &model)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	if req.MsgType == constant.WelcomeMsg {
		caches.WelcomeMessageCacheInstance.DeleteCache(ctx, cardId)
	}

	if req.Id > 0 {
		kflog.AddKFLog(cardId, req.MsgType, "创建了话术", utils.ClientIP(ctx))
	} else {
		kflog.AddKFLog(cardId, req.MsgType, "更新了话术", utils.ClientIP(ctx))
	}

	c.Success(ctx, nil)
}

func (c *WelcomeMsgController) ListAll(ctx *gin.Context) {
	var req kfbackend.ListAllRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()
	cardId := common.GetKFCardID(reqCtx)
	repo := repository.KfWelcomeMessageRepository{}

	data, cnt, err := repo.List(reqCtx, cardId, req.MsgType, req.GetPage(), req.GetPageSize())
	if err != nil {
		c.Error(ctx, err)
		return
	}

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
	cardId := common.GetKFCardID(reqCtx)
	repo := repository.KfWelcomeMessageRepository{}

	deleted, err := repo.Delete(reqCtx, cardId, req.Id)
	if err != nil {
		return
	}

	caches.WelcomeMessageCacheInstance.DeleteCache(ctx, cardId)
	kflog.AddKFLog(cardId, deleted.MsgType, "删除了话术", utils.ClientIP(ctx))

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

	var msgRepo repository.KfWelcomeMessageRepository

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

	caches.WelcomeMessageCacheInstance.DeleteCache(ctx, cardId)
	kflog.AddKFLog(cardId, "话术复制", "从卡密:"+kflog.MaskContent(req.CardID)+" 复制了话术", utils.ClientIP(ctx))
	tx.Commit()
	c.Success(ctx, nil)
}
