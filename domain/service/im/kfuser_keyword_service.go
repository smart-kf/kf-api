package im

import (
	"context"
	"encoding/json"
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

// 智能回复处理.
type kfUserKeywordService struct {
	KfImBaseService
	keywordData dto.KeywordData

	kfMsg   *dao.KFMessage
	userMsg *dao.KFMessage
}

func (s *kfUserKeywordService) Init(ctx context.Context, msg *dto.Message) error {
	s.guestId = msg.Token
	if err := s.setCardIdByGuestId(ctx); err != nil {
		return err
	}
	var keywordData dto.KeywordData
	err := json.Unmarshal([]byte(msg.Content), &keywordData)
	if err != nil {
		return err
	}
	s.keywordData = keywordData

	if err := s.setKfUserSessionId(ctx); err != nil {
		return err
	}
	return nil
}

// deleteSession 删除前台用户的session，并且设置离线
func (s *kfUserKeywordService) deleteSession(ctx context.Context) {
	// 删除用户session.
	caches.ImSessionCacheInstance.DeleteKffeOnlineSession(ctx, s.cardId, s.guestId)
	// 修改客户在线状态: set
	_ = caches.UserOnLineCacheInstance.SetUserOffline(ctx, s.cardId, s.guestId)
}

// Handle 推送给客服后台，客户离线消息
func (s *kfUserKeywordService) Handle(ctx context.Context) {
	if s.keywordData.Id == 0 {
		return
	}

	// 查询消息.
	msg := caches.WelcomeMessageCacheInstance.FindSmartMsg(ctx, s.cardId, int64(s.keywordData.Id))
	if msg == nil {
		return
	}

	if err := s.saveMessage(msg.Type, msg.Keyword, msg.Content); err != nil {
		return
	}

	// 2. 更新用户信息
	if err := s.updateUserLastMessage(ctx, s.kfMsg.MsgId); err != nil {
		return
	}

	// 推送给前台
	pushMessage := dto.NewReplyMessage(
		constant.PlatformKfFe,
		string(s.kfMsg.MsgType),
		s.kfMsg.MsgId,
		s.kfMsg.Content,
		s.guestId,
	)
	s.pushMessage(ctx, constant.EventMessage, pushMessage, s.guestSessionId)

	sessionIds := s.getKfbeSessionIds(ctx, s.cardId)
	if len(sessionIds) == 0 {
		return
	}
	pushMessage.Platform = constant.PlatformKfBe
	pushMessage.Content = s.userMsg.Content
	pushMessage.MsgId = s.userMsg.MsgId
	pushMessage.MsgType = string(s.userMsg.MsgType)
	s.pushMessage(ctx, constant.EventMessage, pushMessage, sessionIds...)

	time.Sleep(1 * time.Second)
	pushMessageBackend := dto.NewReplyMessage(
		constant.PlatformKfBe,
		string(s.kfMsg.MsgType),
		s.kfMsg.MsgId,
		s.kfMsg.Content,
		s.guestId,
	)

	s.pushMessage(ctx, constant.EventMessage, pushMessageBackend, sessionIds...)
}

func (s *kfUserKeywordService) saveMessage(typ string, keyword string, content string) error {
	var msgRepository = repository.KFMessageRepository{}
	// 插入db.
	kfUserMsg := &dao.KFMessage{
		MsgId:   uuid2.NewString(),
		MsgType: common.MessageType(typ),
		GuestId: s.guestId,
		Content: keyword,
		CardId:  s.cardId,
		IsKf:    constant.IsNotKf,
	}
	if err := msgRepository.SaveOne(context.Background(), kfUserMsg); err != nil {
		xlogger.Error(context.Background(), "msgHandler 插入消息失败", xlogger.Err(err))
		return err
	}
	// 插入db.
	kfMsg := &dao.KFMessage{
		MsgId:   uuid2.NewString(),
		MsgType: common.MessageType(typ),
		GuestId: s.guestId,
		Content: content,
		CardId:  s.cardId,
		IsKf:    constant.IsKf,
	}
	if err := msgRepository.SaveOne(context.Background(), kfMsg); err != nil {
		xlogger.Error(context.Background(), "msgHandler 插入消息失败", xlogger.Err(err))
		return err
	}

	s.userMsg = kfUserMsg
	s.kfMsg = kfMsg
	return nil
}
