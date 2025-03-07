package im

import (
	"context"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/dto"
)

type KfOfflineService struct {
	msg *dto.Message
}

func (s *KfOfflineService) Init(ctx context.Context, msg *dto.Message) error {
	s.msg = msg
	return nil
}

func (s *KfOfflineService) Handle(ctx context.Context) {
	caches.ImSessionCacheInstance.DeleteKfBeOnlineSession(ctx, s.msg.Token, s.msg.SessionId)
}
