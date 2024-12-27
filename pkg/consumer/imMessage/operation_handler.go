package imMessage

import (
	"context"
	xlogger "github.com/clearcodecn/log"
	"github.com/smart-fm/kf-api/pkg/goim/api"
)

var operationHandlers = make(map[int32]func(ctx context.Context, op int32, v []byte) error)

func init() {
	operationHandlers[api.OpClientMsg] = handleClientMsg
}

func handleOperation(ctx context.Context, op int32, body []byte) error {
	h, ok := operationHandlers[op]
	if !ok {
		return nil
	}
	return h(ctx, op, body)
}

func handleClientMsg(ctx context.Context, op int32, v []byte) error {
	xlogger.Info(ctx, "收到客户端消息", xlogger.Any("msg", string(v)))
	msg, err := api.FromBytes(v)
	if err != nil {
		xlogger.Error(ctx, "api.FromBytes failed: %+v", xlogger.Any("msg", string(v)))
		return err
	}
	msgService := NewMessageService()
	switch msg.Type {
	case api.MessageTypeDefineMsg:
		m, err := msg.GetMsg()
		if err != nil {
			return err
		}
		err = msgService.ReceiveMessage(ctx, op, m)
	case api.MessageTypeDefineRead:
		m, err := msg.GetRead()
		if err != nil {
			return err
		}
		_ = m
	case api.MessageTypeDefineOnline:
		m, err := msg.GetOnline()
		if err != nil {
			return err
		}
		_ = m
	case api.MessageTypeDefineOffline:
		m, err := msg.GetOffline()
		if err != nil {
			return err
		}
		_ = m
	}

	return nil
}
