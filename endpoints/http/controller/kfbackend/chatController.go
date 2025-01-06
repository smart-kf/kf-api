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

	var repo repository.KFExternalUserRepository

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

	msgIDs := lo.Map(extUsers, func(item *dao.KFExternalUser, index int) uint64 {
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

	chats := lo.Map(extUsers, func(item *dao.KFExternalUser, index int) kfbackend.Chat {
		return extUser2ChatVO(item, lastMsgMap)
	})

	c.Success(ctx, kfbackend.ChatListResponse{
		Chats: chats,
	})
}

func extUser2ChatVO(u *dao.KFExternalUser, lastMsgMap map[uint64]*dao.KFMessage) kfbackend.Chat {
	chat := kfbackend.Chat{
		Type: kfbackend.ChatTypeSingle,
		ExternalUser: kfbackend.ExternalUser{
			Avatar:   u.Avatar,
			NickName: u.NickName,
			IsOnline: false, // TODO 从在离线状态的redis中实时获取
		},
		LastChatAt:   u.LastChatAt,
		UnreadMsgCnt: u.UnreadMsgCnt,
	}

	msg, ok := lastMsgMap[u.LastMsgID]
	if ok {
		chat.LastMessage = kfbackend.Message{
			ID:       msg.ID,
			Content:  msg.Content,
			From:     msg.From,
			FromType: msg.FromType,
			To:       msg.To,
			ToType:   msg.ToType,
			CreateAt: msg.CreatedAt.Unix(),
		}
	}

	return chat
}
