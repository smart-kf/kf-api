package kfbackend

import (
	"context"
	"sort"
	"time"

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

	var repo repository.KFUserRepository

	// in 从redis中获取到的未读访客ids
	var unreadUUIDs []string
	if req.ListType == kfbackend.ChatListTypeUnread {
		uid2Cnt, err := caches.UserUnReadCacheInstance.GetUnReadUsers(ctx, cardID)
		if err != nil {
			xlogger.Error(reqCtx, "查询粉丝未读消息失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
			c.Error(ctx, err)
			return
		}

		uids := lo.MapToSlice(
			uid2Cnt, func(k string, v int64) string {
				return k
			},
		)
		if len(uids) == 0 {
			c.Success(
				ctx, kfbackend.ChatListResponse{},
			)
			return
		}
	}

	users, err := repo.List(
		reqCtx, &repository.ListUserOption{
			CardID:      cardID,
			SearchBy:    req.SearchBy,
			UnreadUUIDs: unreadUUIDs,
			ListType:    req.ListType,
			ScrollRequest: &common.ScrollRequest{
				Sorters: []common.Sorter{
					{
						Key: "top_at", // 1. 置顶时间
						Asc: false,
					},
					{
						Key: "last_chat_at", // 2. 最近聊天时间
						Asc: false,
					},
				},
				ScrollID: req.ScrollID,
				PageSize: req.PageSize,
			},
		},
	)
	if err != nil {
		xlogger.Error(reqCtx, "查询粉丝失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	msgIDs := lo.Map(
		users, func(item *dao.KfUser, index int) uint64 {
			return item.LastMsgID
		},
	)

	var msgRepo repository.KFMessageRepository
	msgs, err := msgRepo.ByIDs(reqCtx, cardID, msgIDs...)
	if err != nil {
		xlogger.Error(
			reqCtx,
			"查询最近消息失败",
			xlogger.Err(err),
			xlogger.Any("cardId", cardID),
			xlogger.Any("msgIDs", msgIDs),
		)
	}

	lastMsgMap := lo.SliceToMap(
		msgs, func(item *dao.KFMessage) (uint64, *dao.KFMessage) {
			return uint64(item.ID), item
		},
	)

	uids := make([]string, 0, len(users))

	chats := lo.Map(
		users, func(item *dao.KfUser, index int) *kfbackend.Chat {
			uids = append(uids, item.UUID)
			return user2ChatVO(reqCtx, item, lastMsgMap)
		},
	)

	// 在线状态.
	if len(uids) > 0 {
		onlineMap, err := caches.UserOnLineCacheInstance.IsUsersOnline(reqCtx, cardID, uids)
		if err == nil {
			lo.ForEach(
				chats, func(item *kfbackend.Chat, index int) {
					item.User.IsOnline = onlineMap[item.User.UUID]
				},
			)
		}
	}

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

//  websocket 来做.
// func (c *ChatController) MsgsRead(ctx *gin.Context) {
// 	var req kfbackend.ReadMsgRequest
// 	if !c.BindAndValidate(ctx, &req) {
// 		return
// 	}
//
// 	if len(req.MsgIDs) == 0 {
// 		// do nothing
// 		c.Success(ctx, kfbackend.ReadMsgResponse{})
// 		return
// 	}
//
// 	reqCtx := ctx.Request.Context()
// 	cardID := common.GetKFCardID(ctx)
//
// 	var repo repository.KFMessageRepository
//
// 	msgs, err := repo.ByIDs(reqCtx, cardID, req.MsgIDs...)
// 	if err != nil {
// 		xlogger.Error(
// 			reqCtx,
// 			"ByIDs",
// 			xlogger.Err(err),
// 			xlogger.Any("cardId", cardID),
// 			xlogger.Any("ids", req.MsgIDs),
// 		)
// 		c.Error(ctx, err)
// 		return
// 	}
//
// 	userIDs := lo.FilterMap(
// 		msgs, func(item *dao.KFMessage, index int) (string, bool) {
// 			if item.FromType == dao.ChatObjTypeUser && len(item.From) > 0 {
// 				return item.From, true
// 			}
// 			return "", false
// 		},
// 	)
//
// 	userIDs = lo.Uniq(userIDs)
//
// 	eg := errgroup.Group{}
// 	eg.SetLimit(len(userIDs) / 3)
// 	for _, userID := range userIDs {
// 		eg.Go(
// 			func() error {
// 				return caches.UserUnReadCacheInstance.IncrUserUnRead(reqCtx, cardID, userID, -1)
// 			},
// 		)
// 	}
//
// 	if err := eg.Wait(); err != nil {
// 		xlogger.Error(
// 			reqCtx,
// 			"删除已读失败",
// 			xlogger.Err(err),
// 			xlogger.Any("cardId", cardID),
// 			xlogger.Any("ids", req.MsgIDs),
// 		)
// 		c.Error(ctx, err)
// 		return
// 	}
//
// 	err = repo.BatchUpdateReadAt(reqCtx, req.MsgIDs, time.Now().Unix())
// 	if err != nil {
// 		xlogger.Error(
// 			reqCtx,
// 			"更新已读时间失败",
// 			xlogger.Err(err),
// 			xlogger.Any("cardId", cardID),
// 			xlogger.Any("ids", req.MsgIDs),
// 		)
// 		return
// 	}
//
// 	c.Success(ctx, kfbackend.ReadMsgResponse{})
// }

func user2ChatVO(ctx context.Context, u *dao.KfUser, lastMsgMap map[uint64]*dao.KFMessage) *kfbackend.Chat {
	unreadCnt, err := caches.UserUnReadCacheInstance.GetUserUnRead(ctx, u.CardID, u.UUID)
	if err != nil {
		xlogger.Error(ctx, "GetUserUnRead err", xlogger.Err(err))
	}
	chat := kfbackend.Chat{
		User:         user2VO(ctx, u),
		LastChatAt:   u.LastChatAt,
		UnreadMsgCnt: unreadCnt,
	}

	msg, ok := lastMsgMap[u.LastMsgID]
	if ok {
		chat.LastMessage = msg2VO(msg)
	}

	return &chat
}

func user2VO(ctx context.Context, u *dao.KfUser) kfbackend.User {
	vo := kfbackend.User{}
	copier.Copy(&vo, u)
	return vo
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
