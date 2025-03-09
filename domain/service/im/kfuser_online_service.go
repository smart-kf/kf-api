package im

import (
	"context"

	xlogger "github.com/clearcodecn/log"
	uuid2 "github.com/google/uuid"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type KFUserOnlineService struct {
	KfImBaseService

	msg *dto.Message
}

func (s *KFUserOnlineService) Init(ctx context.Context, msg *dto.Message) error {
	s.guestId = msg.Token
	s.guestSessionId = msg.SessionId
	s.msg = msg

	err := s.setCardIdByGuestId(ctx)
	if err != nil {
		xlogger.Error(ctx, "处理消息失败-KFUserOnlineService-Init", xlogger.Any("msg", msg))
		return err
	}
	return nil
}

func (s *KFUserOnlineService) saveSession(ctx context.Context) {
	// 将客户session存储至redis
	caches.ImSessionCacheInstance.SetKffeOnlineSession(ctx, s.cardId, s.guestId, s.guestSessionId)
}

func (s *KFUserOnlineService) setUserOnline(ctx context.Context) {
	// 修改客户在线状态: set
	_ = caches.UserOnLineCacheInstance.SetUserOnline(ctx, s.cardId, s.guestId)
}

func (s *KFUserOnlineService) Handle(ctx context.Context) {
	s.saveSession(ctx)
	s.setUserOnline(ctx)

	kfSessionIds := s.getKfbeSessionIds(ctx, s.cardId)
	if len(kfSessionIds) == 0 {
		return
	}
	// 查找前台用户信息.
	dbUser, err := caches.KfUserCacheInstance.GetDBUser(ctx, s.cardId, s.guestId)
	if err != nil {
		xlogger.Error(ctx, "处理消息失败-GetDBUser", xlogger.Err(err))
		return
	}
	onlineMsg := dto.NewGuestOnlineMessage(dbUser.NickName, dbUser.UUID, dbUser.Avatar, s.cardId, s.guestSessionId)
	s.pushMessage(ctx, constant.EventOnline, onlineMsg, kfSessionIds...)

	s.pushWelcomeMessage(ctx, dbUser)
	return
}

func (s *KFUserOnlineService) pushWelcomeMessage(ctx context.Context, dbUser *dao.KfUser) {
	ok := caches.WelcomeMessageCacheInstance.GetUserNeedSendMsg(ctx, s.cardId, dbUser.UUID)
	if !ok {
		return
	}
	caches.WelcomeMessageCacheInstance.DelUserNeedSendMsg(ctx, s.cardId, dbUser.UUID)
	// 查找欢迎语.
	welcomeMsg := caches.WelcomeMessageCacheInstance.FindWelcomeMessages(ctx, s.cardId)
	if len(welcomeMsg) == 0 {
		return
	}
	dbMsgs := make([]*dao.KFMessage, 0, len(welcomeMsg))
	for _, msg := range welcomeMsg {
		dbMsg := dao.KFMessage{
			MsgId:   uuid2.NewString(),
			MsgType: common.MessageType(msg.Type),
			GuestId: dbUser.UUID,
			CardId:  s.cardId,
			Content: msg.Content,
			IsKf:    constant.IsKf,
		}
		dbMsgs = append(dbMsgs, &dbMsg)
	}

	var msgRepo repository.KFMessageRepository
	if err := msgRepo.BatchCreate(ctx, dbMsgs); err != nil {
		xlogger.Error(ctx, "处理消息失败-BatchCreate", xlogger.Err(err))
		return
	}
	for _, msg := range dbMsgs {
		pushMsg := dto.NewPushMessage(string(msg.MsgType), msg.MsgId, msg.Content, dbUser)
		s.pushMessage(ctx, constant.EventMessage, pushMsg, s.guestSessionId)
	}
	s.updateUserLastMessage(ctx, dbMsgs[len(dbMsgs)-1].MsgId)
}
