package im

import (
	"context"
	"time"

	xlogger "github.com/clearcodecn/log"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/pkg/wsclient"
)

type KfImBaseService struct {
	cardId         string
	guestId        string
	kfSessionId    string
	guestSessionId string
}

func (s *KfImBaseService) setCardIdByGuestId(ctx context.Context) error {
	cardId, err := caches.ImSessionCacheInstance.GetCardIDByKFFEToken(ctx, s.guestId)
	if err != nil {
		xlogger.Error(ctx, "处理消息失败-GetCardIDByToken", xlogger.Err(err))
		return err
	}
	s.cardId = cardId
	return nil
}

func (s *KfImBaseService) setCardIdByBackendToken(ctx context.Context, backendToken string) error {
	cardId, err := caches.KfAuthCacheInstance.GetBackendToken(ctx, backendToken)
	if err != nil {
		xlogger.Error(ctx, "处理消息失败-GetBackendToken", xlogger.Err(err))
		return err
	}
	s.cardId = cardId
	return nil
}

func (s *KfImBaseService) getKfbeSessionIds(ctx context.Context, cardId string) []string {
	sessionIds, err := caches.ImSessionCacheInstance.GetKfbeSessionIds(ctx, cardId)
	if err != nil {
		xlogger.Error(ctx, "处理消息失败-getKfbeSessionIds", xlogger.Err(err))
		return nil
	}
	if len(sessionIds) == 0 {
		return nil
	}
	return sessionIds
}

func (s *KfImBaseService) setKfUserSessionId(ctx context.Context) error {
	sid, err := caches.ImSessionCacheInstance.GetKFFESessionIdByUserId(ctx, s.cardId, s.guestId)
	if err != nil {
		xlogger.Error(ctx, "处理消息失败-getKfbeSessionIds", xlogger.Err(err))
		return err
	}
	s.guestSessionId = sid
	return nil
}

func (s *KfImBaseService) pushMessage(
	ctx context.Context,
	event string,
	message *dto.Message,
	sessionIds ...string,
) {
	var pushClient = wsclient.WsClient{}
	if err := pushClient.Push(ctx, event, message, sessionIds...); err != nil {
		xlogger.Error(ctx, "pushManyMessage-failed", xlogger.Err(err))
	}
}

func (s *KfImBaseService) updateUserLastMessage(ctx context.Context, msgId string) error {
	userRepository := repository.KFUserRepository{}
	// 更新用户信息.
	opt := repository.UpdateColOption{}
	opt.UUID = s.guestId
	opt.LastChatAt = time.Now().Unix()
	opt.LastMessageId = msgId
	err := userRepository.UpdateCol(ctx, s.cardId, opt)
	if err != nil {
		xlogger.Error(ctx, "msgHandler 更新用户信息失败", xlogger.Err(err))
		return err
	}
	return nil
}
