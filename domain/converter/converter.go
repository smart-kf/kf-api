package converter

import (
	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

func DtoMessageToModelMessage(message dto.Message) *dao.KFMessage {
	msg := &dao.KFMessage{
		MsgId:   message.MsgId,
		MsgType: dao.MessageType(message.MsgType),
		GuestId: message.GuestId,
		CardId:  message.KfId,
		Content: message.Content,
		IsKf:    message.IsKf,
	}
	return msg
}
