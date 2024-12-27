package imMessage

import (
	"context"
	"encoding/json"
	xlogger "github.com/clearcodecn/log"
	"github.com/smart-fm/kf-api/pkg/db"
	"github.com/smart-fm/kf-api/pkg/goim/api"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) ReceiveMessage(ctx context.Context, op int32, msg *api.MessageTypeMsg) error {
	// 1. 保存到数据库.
	tx := db.GetDBFromContext(ctx)
	kfMessage := &db.KFMessage{
		MsgType: msg.MsgType,
		KfId:    msg.KfId,
		GuestId: msg.GuestId,
		Content: msg.Content,
		IsKF:    msg.IsKF,
	}

	if err := tx.Create(kfMessage).Error; err != nil {
		xlogger.Error(ctx, "Create KFMessage failed", xlogger.Err(err))
		return err
	}

	// push.
	logicClient := api.GetLogicClient()

	var pushMid int64

	// 2. 推送消息给对方.
	if api.IsKFMessage(msg.IsKF) {
		// 推送给客户.
		if msg.GuestId == 0 { // 找不到，不推送.
			return nil
		}
		pushMid = msg.GuestId
	} else {
		// 推送给客服
		if msg.KfId == 0 { // 找不到，不推送.
			return nil
		}
		pushMid = msg.KfId
	}
	data, _ := json.Marshal(msg)
	raw := json.RawMessage(data)
	err := logicClient.PushMids(ctx, op, []int64{pushMid}, &api.Message{
		Type: api.MessageTypeDefineMsg,
		Data: &raw,
	})
	if err != nil {
		xlogger.Error(ctx, "ReceiveMessage-pushMid failed", xlogger.Err(err))
	}
	return nil
}
