package im

import (
	"context"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/dto"
)

type KfReadService struct {
	KfImBaseService
}

func (s *KfReadService) Init(ctx context.Context, msg *dto.Message) error {
	if err := s.setCardIdByBackendToken(ctx, msg.Token); err != nil {
		return err
	}
	s.guestId = msg.GuestId
	return nil
}

func (s *KfReadService) clearReadMessages(ctx context.Context) {
	caches.UserUnReadCacheInstance.IncrUserUnRead(ctx, s.cardId, s.guestId, -1)
}

func (s *KfReadService) Handle(ctx context.Context) {
	s.clearReadMessages(ctx)
	return
}
