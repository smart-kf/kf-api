package im

import (
	"context"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
)

type KfUserOfflineService struct {
	KfImBaseService
}

func (s *KfUserOfflineService) Init(ctx context.Context, msg *dto.Message) error {
	s.guestId = msg.Token
	if err := s.setCardIdByGuestId(ctx); err != nil {
		return err
	}
	return nil
}

// deleteSession 删除前台用户的session，并且设置离线
func (s *KfUserOfflineService) deleteSession(ctx context.Context) {
	// 删除用户session.
	caches.ImSessionCacheInstance.DeleteKffeOnlineSession(ctx, s.cardId, s.guestId)
	// 修改客户在线状态: set
	_ = caches.UserOnLineCacheInstance.SetUserOffline(ctx, s.cardId, s.guestId)
}

// Handle 推送给客服后台，客户离线消息
func (s *KfUserOfflineService) Handle(ctx context.Context) {
	s.deleteSession(ctx)
	// 推送离线事件.
	sessionIds := s.getKfbeSessionIds(ctx, s.cardId)
	if len(sessionIds) == 0 {
		return
	}
	msg := dto.NewGuestOfflineMessage(s.guestId, s.cardId, s.guestSessionId)
	s.pushMessage(ctx, constant.EventOffline, msg, sessionIds...)
}
