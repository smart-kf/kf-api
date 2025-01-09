package kfbackend

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/middleware"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"time"
)

type ChatController struct {
	BaseController
}

func (c *ChatController) List(ctx *gin.Context) {
	var req kfbackend.ChatListRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()
	cardID := middleware.GetKFCardID(ctx)

	var repo repository.KFUserRepository

	extUsers, err := repo.List(reqCtx, &repository.ListExtUserOption{
		CardID:   cardID,
		SearchBy: req.SearchBy,
		ListType: req.ListType,
		ScrollRequest: &common.ScrollRequest{
			Key:      "last_chat_at",
			Asc:      false,
			ScrollID: req.ScrollID,
			PageSize: req.PageSize,
		},
	})
	if err != nil {
		xlogger.Error(reqCtx, "查询客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	msgIDs := lo.Map(extUsers, func(item *dao.KfUser, index int) uint64 {
		return item.LastMsgID
	})

	var msgRepo repository.KFMessageRepository
	msgs, err := msgRepo.ByIDs(reqCtx, cardID, msgIDs...)
	if err != nil {
		xlogger.Error(reqCtx, "查询最近消息失败", xlogger.Err(err), xlogger.Any("cardId", cardID), xlogger.Any("msgIDs", msgIDs))
	}

	lastMsgMap := lo.SliceToMap(msgs, func(item *dao.KFMessage) (uint64, *dao.KFMessage) {
		return item.ID, item
	})

	chats := lo.Map(extUsers, func(item *dao.KfUser, index int) kfbackend.Chat {
		return extUser2ChatVO(item, lastMsgMap)
	})

	c.Success(ctx, kfbackend.ChatListResponse{
		Chats: chats,
	})
}

func (c *ChatController) Msgs(ctx *gin.Context) {
	var req kfbackend.MsgListRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	if len(req.FromTos) == 0 {
		c.Success(ctx, kfbackend.MsgListResponse{})
		return
	}

	reqCtx := ctx.Request.Context()
	cardID := middleware.GetKFCardID(ctx)

	var repo repository.KFMessageRepository

	extUsers, err := repo.List(reqCtx, &repository.ListMsgOption{
		CardID:  cardID,
		FromTos: req.FromTos,
		ScrollRequest: &common.ScrollRequest{
			Key:      "id",
			Asc:      req.Asc,
			ScrollID: req.ScrollID,
			PageSize: req.PageSize,
		},
	})
	if err != nil {
		xlogger.Error(reqCtx, "查询客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	msgs := lo.Map(extUsers, func(item *dao.KFMessage, index int) kfbackend.Message {
		return msg2VO(item)
	})

	c.Success(ctx, kfbackend.MsgListResponse{
		Messages: msgs,
	})
}

func (c *ChatController) MsgsRead(ctx *gin.Context) {
	var req kfbackend.ReadMsgRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	if len(req.MsgIDs) == 0 {
		// do nothing
		c.Success(ctx, kfbackend.ReadMsgResponse{})
		return
	}

	reqCtx := ctx.Request.Context()
	cardID := middleware.GetKFCardID(ctx)

	var repo repository.KFMessageRepository

	err := repo.BatchUpdateReadAt(reqCtx, req.MsgIDs, time.Now().Unix())
	if err != nil {
		xlogger.Error(reqCtx, "更新已读时间失败", xlogger.Err(err), xlogger.Any("cardId", cardID), xlogger.Any("ids", req.MsgIDs))
		c.Error(ctx, err)
		return
	}

	c.Success(ctx, kfbackend.ReadMsgResponse{})
}

func (c *ChatController) UserOp(ctx *gin.Context) {
	var req kfbackend.BatchOpUserRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	if len(req.UserIDs) == 0 {
		// do nothing
		c.Success(ctx, kfbackend.BatchOpUserResponse{})
		return
	}

	reqCtx := ctx.Request.Context()
	cardID := middleware.GetKFCardID(ctx)

	var repo repository.KFUserRepository

	u := dao.KfUser{}
	switch req.Op {
	case kfbackend.UserOpTop:
		u.TopAt = time.Now().Unix()
	case kfbackend.UserOpTopUndo:
		u.TopAt = 0
	case kfbackend.UserOpBlock:
		u.BlockAt = time.Now().Unix()
	case kfbackend.UserOpBlockUndo:
		u.BlockAt = 0
	}

	err := repo.BatchUpdate(reqCtx, req.UserIDs, u)
	if err != nil {
		xlogger.Error(reqCtx, "更新访客失败", xlogger.Err(err), xlogger.Any("cardId", cardID), xlogger.Any("ids", req.UserIDs))
		c.Error(ctx, err)
		return
	}

	c.Success(ctx, kfbackend.BatchOpUserResponse{})
}

func (c *ChatController) UserUpdate(ctx *gin.Context) {
	var req kfbackend.UpdateUserRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	if req.ID <= 0 {
		// do nothing
		c.Success(ctx, kfbackend.UpdateUserResponse{})
		return
	}

	reqCtx := ctx.Request.Context()
	cardID := middleware.GetKFCardID(ctx)

	var repo repository.KFUserRepository

	u := dao.KfUser{
		RemarkName: req.RemarkName,
		Mobile:     req.Mobile,
		Comments:   req.Comments,
	}

	err := repo.BatchUpdate(reqCtx, []uint{req.ID}, u)
	if err != nil {
		xlogger.Error(reqCtx, "更新访客失败", xlogger.Err(err), xlogger.Any("cardId", cardID), xlogger.Any("ids", req.ID))
		c.Error(ctx, err)
		return
	}

	c.Success(ctx, kfbackend.BatchOpUserResponse{})
}

func extUser2ChatVO(u *dao.KfUser, lastMsgMap map[uint64]*dao.KFMessage) kfbackend.Chat {
	chat := kfbackend.Chat{
		Type: kfbackend.ChatTypeSingle,
		User: kfbackend.User{
			Avatar:   u.Avatar,
			NickName: u.NickName,
			IsOnline: false, // TODO 从在离线状态的redis中实时获取
		},
		LastChatAt:   u.LastChatAt,
		UnreadMsgCnt: u.UnreadMsgCnt,
	}

	msg, ok := lastMsgMap[u.LastMsgID]
	if ok {
		chat.LastMessage = msg2VO(msg)
	}

	return chat
}

func msg2VO(m *dao.KFMessage) kfbackend.Message {
	vo := kfbackend.Message{
		ID:       m.ID,
		Content:  m.Content,
		From:     m.From,
		FromType: m.FromType,
		To:       m.To,
		ToType:   m.ToType,
		ReadAt:   m.ReadAt,
		CreateAt: m.CreatedAt.Unix(),
	}

	return vo
}
