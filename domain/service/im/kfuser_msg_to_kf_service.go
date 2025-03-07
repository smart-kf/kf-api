package im

import (
	"context"
	"time"

	xlogger "github.com/clearcodecn/log"
	uuid2 "github.com/google/uuid"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type KfUserMsgToKfService struct {
	KfImBaseService
	msg        *dto.Message
	newMessage *dao.KFMessage
}

func (s *KfUserMsgToKfService) Init(ctx context.Context, msg *dto.Message) error {
	s.guestId = msg.Token
	if err := s.setCardIdByGuestId(ctx); err != nil {
		return err
	}
	s.msg = msg
	return nil
}

func (s *KfUserMsgToKfService) saveMessage() error {
	var msgRepository = repository.KFMessageRepository{}
	// 插入db.
	newMessage := &dao.KFMessage{
		MsgId:   uuid2.NewString(),
		MsgType: common.MessageType(s.msg.MsgType),
		GuestId: s.guestId,
		Content: s.msg.Content,
		CardId:  s.cardId,
		IsKf:    constant.IsNotKf,
	}
	if err := msgRepository.SaveOne(context.Background(), newMessage); err != nil {
		xlogger.Error(context.Background(), "msgHandler 插入消息失败", xlogger.Err(err))
		return err
	}
	s.newMessage = newMessage
	return nil
}

func (s *KfUserMsgToKfService) updateUserInfo(ctx context.Context) error {
	userRepository := repository.KFUserRepository{}
	// 更新用户信息.
	opt := repository.UpdateColOption{}
	opt.UUID = s.guestId
	opt.LastChatAt = time.Now().Unix()
	opt.LastMessageId = s.newMessage.MsgId
	err := userRepository.UpdateCol(ctx, s.cardId, opt)
	if err != nil {
		xlogger.Error(ctx, "msgHandler 更新用户信息失败", xlogger.Err(err))
		return err
	}
	return nil
}

func (s *KfUserMsgToKfService) increaseUnRead(ctx context.Context) {
	// 未读 + 1
	caches.UserUnReadCacheInstance.IncrUserUnRead(ctx, s.cardId, s.guestId, 1)
}

// Handle 前台消息发给后台
func (s *KfUserMsgToKfService) Handle(ctx context.Context) {
	// 1. 存储消息入库
	if err := s.saveMessage(); err != nil {
		return
	}
	// 2. 更新用户信息
	if err := s.updateUserInfo(ctx); err != nil {
		return
	}
	s.increaseUnRead(ctx)

	sessionIds := s.getKfbeSessionIds(ctx, s.cardId)
	if len(sessionIds) == 0 {
		s.fallbackToAIAnswer(ctx)
		return
	}
	pushMessage := dto.NewMessage(s.msg, constant.PlatformKfBe)
	s.pushMessage(ctx, constant.EventMessage, pushMessage, sessionIds...)
}

// 智能答复.
func (s *KfUserMsgToKfService) fallbackToAIAnswer(ctx context.Context) {
	// TODO:: implement me.
}
