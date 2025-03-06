package kfbackend

import (
	"sort"
	"time"
	"unicode/utf8"

	uuid2 "github.com/google/uuid"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/endpoints/cron/kflog"
	"github.com/smart-fm/kf-api/infrastructure/httpClient/socketserver"
	"github.com/smart-fm/kf-api/pkg/utils"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
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
	cardID := common.GetKFCardID(ctx)

	var (
		repo   repository.KFUserRepository
		err    error
		helper caches.UnReadHelper
	)

	// in 从redis中获取到的未读访客ids
	var (
		unreadUUIDs   []string
		unreadUserMap map[string]int64
	)
	if req.ListType == kfbackend.ChatListTypeUnread {
		unreadUserMap, unreadUUIDs, err = helper.GetUnReadUserIDs(reqCtx, cardID)
		if err != nil {
			xlogger.Error(ctx, "GetUnReadUserIDs failed", xlogger.Err(err))
		}
	}
	listOption := &repository.ListUserOption{
		CardID:      cardID,
		SearchBy:    req.SearchBy,
		UnreadUUIDs: unreadUUIDs,
		ListType:    req.ListType,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}
	if req.ListType == kfbackend.ChatListTypeBlock {
		listOption.Blocked = true
	}

	users, err := repo.List(
		reqCtx, listOption,
	)
	if err != nil {
		xlogger.Error(reqCtx, "查询粉丝失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	if len(users) == 0 {
		c.Success(ctx, kfbackend.ChatListResponse{})
		return
	}

	var (
		lastMsgMap map[string]*dao.KFMessage
		onlineMap  map[string]bool
		listUids   []string
		lastMsgIds []string
		msgRepo    repository.KFMessageRepository
	)
	lo.ForEach(
		users, func(item *dao.KfUser, index int) {
			if item.LastMessageId != "" {
				lastMsgIds = append(lastMsgIds, item.LastMessageId)
			}
			listUids = append(listUids, item.UUID)
		},
	)
	if len(lastMsgIds) > 0 {
		msgs, err := msgRepo.ByIDs(reqCtx, cardID, lastMsgIds)
		if err != nil {
			xlogger.Error(
				reqCtx,
				"查询最近消息失败",
				xlogger.Err(err),
				xlogger.Any("cardId", cardID),
				xlogger.Any("msgIDs", lastMsgIds),
			)
		}
		lastMsgMap = lo.SliceToMap(
			msgs, func(item *dao.KFMessage) (string, *dao.KFMessage) {
				return item.GuestId, item
			},
		)
	}
	// 在线状态.
	if len(listUids) > 0 {
		// 最近一条消息的用户信息.
		onlineMap, _ = caches.UserOnLineCacheInstance.IsUsersOnline(reqCtx, cardID, listUids)
	}
	var chats []*kfbackend.Chat
	lo.ForEach(
		users, func(item *dao.KfUser, index int) {
			var (
				user kfbackend.User
				chat kfbackend.Chat
			)
			copier.Copy(&user, item)

			if onlineMap != nil {
				user.IsOnline = onlineMap[item.UUID]
			}
			if lastMsgMap != nil {
				msg, ok := lastMsgMap[item.UUID]
				if ok {
					if msg.MsgType == common.MessageTypeText {
						if utf8.RuneCountInString(msg.Content) > 10 {
							runes := []rune(msg.Content)
							chat.LastMessage = string(runes[:10])
						}
					}
					if msg.MsgType == common.MessageTypeVideo {
						chat.LastMessage = "[视频消息]"
					}
					if msg.MsgType == common.MessageTypeVoice {
						chat.LastMessage = "[语音消息]"
					}
				}
			}
			if unreadUserMap != nil {
				unread, ok := unreadUserMap[item.UUID]
				if ok {
					chat.UnreadMsgCnt = unread
				}
			}
			chat.User = user
			chats = append(chats, &chat)
		},
	)

	c.Success(
		ctx, kfbackend.ChatListResponse{
			Chats: chats,
		},
	)
}

func (c *ChatController) Msgs(ctx *gin.Context) {
	var req kfbackend.MsgListRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()
	cardID := common.GetKFCardID(ctx)

	var repo repository.KFMessageRepository

	var ti time.Time
	if req.LastMsgTime != 0 {
		ti = time.Unix(req.LastMsgTime, 0)
	}

	msgsDTO, err := repo.List(
		reqCtx, &repository.ListMsgOption{
			CardID:      cardID,
			GuestId:     req.GuestId,
			LastMsgTime: ti,
			PageSize:    int(req.PageSize),
		},
	)
	if err != nil {
		xlogger.Error(reqCtx, "查询消息失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	sort.Slice(
		msgsDTO, func(i, j int) bool {
			return msgsDTO[i].CreatedAt.Unix() < msgsDTO[j].CreatedAt.Unix()
		},
	)

	msgsVO := lo.Map(
		msgsDTO, func(item *dao.KFMessage, index int) *kfbackend.Message {
			return msg2VO(item)
		},
	)

	// 清空未读消息.
	caches.UserUnReadCacheInstance.IncrUserUnRead(reqCtx, cardID, req.GuestId, -1)

	c.Success(
		ctx, kfbackend.MsgListResponse{
			Messages: msgsVO,
		},
	)
}

func msg2VO(m *dao.KFMessage) *kfbackend.Message {
	vo := &kfbackend.Message{
		MsgId:   m.MsgId,
		MsgType: m.MsgType,
		GuestId: m.GuestId,
		CardId:  m.CardId,
		Content: m.Content,
		IsKf:    m.IsKf,
		MsgTime: m.CreatedAt.Unix(),
	}
	return vo
}

func (c *ChatController) BatchSend(ctx *gin.Context) {
	var req kfbackend.BatchSendRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()

	cardId := common.GetKFCardID(reqCtx)
	sessionIdMap, err := caches.ImSessionCacheInstance.GetKffeSessionIds(ctx, cardId, req.GuestId)
	if err != nil {
		return
	}
	var sessionIds []string
	for _, item := range sessionIdMap {
		sessionIds = append(sessionIds, item)
	}
	// 推给前台
	newMessage := dto.Message{
		MessageBase: dto.MessageBase{
			Event:    constant.EventMessage,
			Platform: constant.PlatformKfFe,
		},
		MsgType: string(req.Message.MsgType),
		MsgId:   req.Message.MsgId,
		GuestId: req.Message.GuestId,
		Content: req.Message.Content,
		IsKf:    constant.IsKf,
	}
	pushMsgRequest := socketserver.PushMessageRequest{
		SessionIds: sessionIds,
		Event:      constant.EventMessage,
	}
	pushMsgRequest.SetData(newMessage)
	if err := socketserver.NewSocketServerClient().PushMessage(reqCtx, &pushMsgRequest); err != nil {
		c.Error(ctx, err)
		return
	}
	// 入库.
	var msgs []*dao.KFMessage
	for _, id := range req.GuestId {
		msgs = append(
			msgs, &dao.KFMessage{
				MsgId:   uuid2.NewString(),
				MsgType: req.Message.MsgType,
				GuestId: id,
				CardId:  cardId,
				Content: req.Message.Content,
				IsKf:    constant.IsKf,
			},
		)
	}

	var repo repository.KFMessageRepository
	if err := repo.BatchCreate(reqCtx, msgs); err != nil {
		c.Error(ctx, err)
		return
	}

	kflog.AddKFLog(cardId, "客户", "群发了一条信息", utils.ClientIP(ctx))

	c.Success(ctx, nil)
}
