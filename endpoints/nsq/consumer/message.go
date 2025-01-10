package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/nsqio/go-nsq"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/dto"
	event2 "github.com/smart-fm/kf-api/domain/event"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/httpClient/socketserver"
)

// MessageConsumer 消息消费者.
type MessageConsumer struct{}

func (m *MessageConsumer) HandleMessage(message *nsq.Message) error {
	fmt.Println("receive a new message --->", string(message.Body))
	// 创建消息，并且给客户端回复消息id.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_ = ctx

	var msg dto.Message
	err := json.Unmarshal(message.Body, &msg)
	if err != nil {
		return err
	}
	switch msg.Event {
	case constant.EventMessage:
		m.handleEventMessage(&msg)
	case constant.EventOnline:
		m.handleOnline(&msg)
	case constant.EventOffline:
		m.handleOffline(&msg)
	}

	return nil
}

func (m *MessageConsumer) handleOnline(msg *dto.Message) {
	ctx := context.Background()
	// 1. 发布 sessionId 事件.
	client := socketserver.NewSocketServerClient()
	{
		var req = socketserver.PushMessageRequest{
			SessionId: msg.SessionId,
			Event:     constant.EventSessionId,
		}
		req.SetData(
			map[string]string{
				"sessionId": msg.SessionId,
			},
		)
		client.PushMessage(context.Background(), &req)
	}
	{
		// 2. 如果是前台用户上线，给后台推送一个上线事件，并且修改redis状态.完了之后，触发event的上线钩子.
		if msg.Platform == constant.PlatformKfFe {
			cardId, err := caches.ImSessionCacheInstance.GetCardIDByToken(ctx, msg.Token)
			if err != nil {
				xlogger.Error(ctx, "处理消息失败-GetCardIDByToken", xlogger.Err(err), xlogger.Any("msg", msg))
				return
			}
			// 将客户session存储至redis
			caches.ImSessionCacheInstance.SetKffeOnlineSession(ctx, cardId, msg.Token, msg.SessionId)
			// 修改客户在线状态: set
			_ = caches.UserOnLineCacheInstance.SetUserOnline(ctx, cardId, msg.Token)

			sessionIds, err := caches.ImSessionCacheInstance.GetKfbeSessionIds(ctx, cardId)
			if err != nil {
				xlogger.Error(ctx, "处理消息失败-GetKfbeSessionIds", xlogger.Err(err), xlogger.Any("msg", msg))
				return
			}
			if len(sessionIds) == 0 {
				return
			}
			// 如果是前台用户上线，给后台推送一个前台用户上线事件.
			var req = socketserver.PushMessageRequest{
				SessionIds: sessionIds,
				Event:      constant.EventOnline,
			}
			// 查找前台用户信息.
			dbUser, err := caches.KfUserCacheInstance.GetDBUser(ctx, cardId, msg.Token)
			if err != nil {
				xlogger.Error(ctx, "处理消息失败-GetDBUser", xlogger.Err(err), xlogger.Any("msg", msg))
				return
			}
			var event = dto.Online{
				MessageBase: dto.MessageBase{
					Event:     "",
					Platform:  "",
					SessionId: "",
					Token:     "",
				},
				GuestName:   dbUser.NickName,
				GuestAvatar: dbUser.Avatar,
				GuestId:     dbUser.UUID,
				IsKf:        constant.IsKf,
				KfId:        cardId,
			}
			req.SetData(event)
			client.PushMessage(ctx, &req)
			// 触发evens
			event2.TriggerEvent(context.Background(), constant.EventOnline, msg.Platform, msg.Token, cardId)
			return
		}
		if msg.Platform == constant.PlatformKfBe {
			// 将后台session存储至redis
			caches.ImSessionCacheInstance.SetKfbeOnlineSession(ctx, msg.Token, msg.SessionId)
			cardId, err := caches.KfAuthCacheInstance.GetBackendToken(ctx, msg.Token)
			if err != nil {
				return
			}
			event2.TriggerEvent(context.Background(), constant.EventOnline, msg.Platform, msg.Token, cardId)
		}
	}
}

// handleOffline 离线事件.
func (m *MessageConsumer) handleOffline(msg *dto.Message) {
	ctx := context.Background()
	// 1. 发布 sessionId 事件.
	client := socketserver.NewSocketServerClient()
	{
		// 2. 如果是前台用户上线，给后台推送一个上线事件，并且修改redis状态.完了之后，触发event的上线钩子.
		if msg.Platform == constant.PlatformKfFe {
			cardId, err := caches.ImSessionCacheInstance.GetCardIDByToken(ctx, msg.Token)
			if err != nil {
				xlogger.Error(ctx, "处理消息失败-GetCardIDByToken", xlogger.Err(err), xlogger.Any("msg", msg))
				return
			}
			// 删除用户session.
			caches.ImSessionCacheInstance.DeleteKffeOnlineSession(ctx, cardId, msg.Token)
			// 修改客户在线状态: set
			_ = caches.UserOnLineCacheInstance.SetUserOffline(ctx, cardId, msg.Token)
			// 推送离线事件.
			sessionIds, err := caches.ImSessionCacheInstance.GetKfbeSessionIds(ctx, cardId)
			if err != nil {
				xlogger.Error(ctx, "处理消息失败-GetKfbeSessionIds", xlogger.Err(err), xlogger.Any("msg", msg))
				return
			}
			if len(sessionIds) == 0 {
				return
			}
			// 如果是前台用户离线，给后台推送一个前台用户离线事件.
			var req = socketserver.PushMessageRequest{
				SessionIds: sessionIds,
				Event:      constant.EventOffline,
			}
			// 查找前台用户信息.
			dbUser, err := caches.KfUserCacheInstance.GetDBUser(ctx, cardId, msg.Token)
			if err != nil {
				xlogger.Error(ctx, "处理消息失败-GetDBUser", xlogger.Err(err), xlogger.Any("msg", msg))
				return
			}
			var event = dto.Online{
				MessageBase: dto.MessageBase{
					Event:     constant.EventOffline,
					Platform:  constant.PlatformKfBe,
					SessionId: "",        // TODO:: implement
					Token:     msg.Token, // TODO:: checking.
				},
				GuestName:   dbUser.NickName,
				GuestAvatar: dbUser.Avatar,
				GuestId:     dbUser.UUID,
				IsKf:        constant.IsKf,
				KfId:        cardId,
			}
			req.SetData(event)
			client.PushMessage(ctx, &req)
			// 触发evens
			event2.TriggerEvent(context.Background(), constant.EventOffline, msg.Platform, msg.Token, cardId)
			return
		}
		if msg.Platform == constant.PlatformKfBe {
			// 将后台session存储至redis
			caches.ImSessionCacheInstance.DeleteKfBeOnlineSession(ctx, msg.Token, msg.SessionId)
			cardId, err := caches.KfAuthCacheInstance.GetBackendToken(ctx, msg.Token)
			if err != nil {
				return
			}
			event2.TriggerEvent(context.Background(), constant.EventOffline, msg.Platform, msg.Token, cardId)
		}
	}
}

// handleEventMessage:: todo.
func (m *MessageConsumer) handleEventMessage(msg *dto.Message) {
	// ctx := context.Background()
	// client := socketserver.NewSocketServerClient()
	// // 1. 回复已收到ACK
	// var req = socketserver.PushMessageRequest{
	// 	SessionId: msg.SessionId,
	// 	Event:     constant.EventMessageAck,
	// }
	// req.SetData(
	// 	&dto.Message{
	// 		MessageBase: dto.MessageBase{
	// 			Event:     constant.EventMessageAck,
	// 			Platform:  msg.Platform,
	// 			SessionId: msg.SessionId,
	// 			Token:     "",
	// 		},
	// 		MsgType:     "",
	// 		MsgId:       "",
	// 		GuestName:   "",
	// 		GuestAvatar: "",
	// 		GuestId:     "",
	// 		Content:     "",
	// 		KfId:        "",
	// 		IsKf:        0,
	// 	},
	// )
	// client.PushMessage(context.Background(), &req)
	//
	// var newMessage Message
	// // 2. 推送给接收方.
	// if msg.Platform == constant.PlatformKfFe {
	// 	// 查询接收方的id.
	// 	var repo repository.KfUserRedisRepository
	// 	user, ok, err := repo.GetDBUser(ctx, msg.Token)
	// 	if err != nil {
	// 		return
	// 	}
	// 	if !ok {
	// 		return
	// 	}
	// 	// 推给后台.
	// 	newMessage = Message{
	// 		Event:       constant.EventMessage,
	// 		Platform:    constant.PlatformKFBackend,
	// 		SessionId:   "",
	// 		Token:       "",
	// 		MsgType:     msg.MsgType,
	// 		MsgId:       msg.MsgId,
	// 		GuestName:   user.NickName,
	// 		GuestAvatar: user.Avatar,
	// 		GuestId:     user.UUID,
	// 		Content:     msg.Content,
	// 		KfId:        "",
	// 		IsKf:        2,
	// 	}
	// } else {
	// 	// 推给前台
	// 	// 推给后台.
	// 	newMessage = Message{
	// 		Event:       constant.EventMessage,
	// 		Platform:    constant.PlatformKFBackend,
	// 		SessionId:   "",
	// 		Token:       "",
	// 		MsgType:     msg.MsgType,
	// 		MsgId:       msg.MsgId,
	// 		GuestName:   user.NickName,
	// 		GuestAvatar: user.Avatar,
	// 		GuestId:     user.UUID,
	// 		Content:     msg.Content,
	// 		KfId:        "",
	// 		IsKf:        2,
	// 	}
	// }
}
