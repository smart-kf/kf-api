package im

import (
	"context"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/dto"
)

type KfOnlineService struct {
	KfImBaseService

	msg *dto.Message
}

func (s *KfOnlineService) Init(ctx context.Context, msg *dto.Message) error {
	if err := s.setCardIdByBackendToken(ctx, msg.Token); err != nil {
		return err
	}
	s.msg = msg
	s.kfSessionId = msg.SessionId
	return nil
}

func (s *KfOnlineService) saveBackendSession(ctx context.Context) {
	// 将后台session存储至redis
	caches.ImSessionCacheInstance.SetKfbeOnlineSession(ctx, s.cardId, s.kfSessionId)
}

func (s *KfOnlineService) Handle(ctx context.Context) {
	s.saveBackendSession(ctx)
	return
}
