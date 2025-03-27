package im

import (
	"context"

	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
)

type MessageHandler interface {
	Init(ctx context.Context, msg *dto.Message) error
	Handle(ctx context.Context)
}

func FactoryMessageHandler(eventType string, platform string, msgType string) MessageHandler {
	switch eventType {
	case constant.EventMessage:
		if msgType == constant.MsgTypeRead {
			return &KfReadService{}
		}
		if msgType == constant.MsgKeyword {
			// return &kfUserKeywordService{}
		}
		if platform == constant.PlatformKfBe {
			return &KfMsgToKfUserService{}
		} else {
			return &KfUserMsgToKfService{}
		}
	case constant.EventOnline:
		if platform == constant.PlatformKfFe {
			return &KFUserOnlineService{}
		} else {
			return &KfOnlineService{}
		}
	case constant.EventOffline:
		if platform == constant.PlatformKfFe {
			return &KfUserOfflineService{}
		} else {
			return &KfOfflineService{}
		}
	}

	return nil
}
