package im

import (
	"context"
	"time"

	xlogger "github.com/clearcodecn/log"
	uuid2 "github.com/google/uuid"

	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type KfMsgToKfUserService struct {
	KfImBaseService
	msg        *dto.Message
	newMessage *dao.KFMessage
}

func (s *KfMsgToKfUserService) Init(ctx context.Context, msg *dto.Message) error {
	s.msg = msg
	s.guestId = msg.GuestId
	if err := s.setCardIdByGuestId(ctx); err != nil {
		return err
	}
	return nil
}

func (s *KfMsgToKfUserService) saveMessage(ctx context.Context) error {
	var msgRepository = repository.KFMessageRepository{}
	// 插入db.
	newMessage := &dao.KFMessage{
		MsgId:   uuid2.NewString(),
		MsgType: common.MessageType(s.msg.MsgType),
		GuestId: s.guestId,
		CardId:  s.cardId,
		Content: s.msg.Content,
		IsKf:    constant.IsKf,
	}
	if err := msgRepository.SaveOne(ctx, newMessage); err != nil {
		xlogger.Error(ctx, "msgHandler 插入消息失败", xlogger.Err(err))
		return err
	}
	s.newMessage = newMessage
	return nil
}

func (s *KfMsgToKfUserService) updateUserInfo(ctx context.Context) error {
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

func (s *KfMsgToKfUserService) Handle(ctx context.Context) {
	if err := s.saveMessage(ctx); err != nil {
		return
	}

	if err := s.updateUserInfo(ctx); err != nil {
		return
	}

	if err := s.setKfUserSessionId(ctx); err != nil {
		return
	}

	if s.guestSessionId == "" {
		return
	}

	newMessage := dto.NewMessage(s.msg, constant.PlatformKfFe)
	s.pushMessage(ctx, constant.EventMessage, newMessage, s.guestSessionId)
}
