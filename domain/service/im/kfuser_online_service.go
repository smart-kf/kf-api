package im

import (
	"context"

	xlogger "github.com/clearcodecn/log"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
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
	return
}
