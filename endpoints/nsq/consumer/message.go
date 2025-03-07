package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"

	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/domain/service/im"
)

// MessageConsumer 消息消费者.
type MessageConsumer struct{}

func (m *MessageConsumer) HandleMessage(message *nsq.Message) error {
	fmt.Println("receive a new message --->", string(message.Body))
	// 创建消息，并且给客户端回复消息id.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_ = ctx

	var msg dto.Message
	err := json.Unmarshal(message.Body, &msg)
	if err != nil {
		return err
	}

	h := im.FactoryMessageHandler(msg.Event, msg.Platform, msg.MsgType)
	if h == nil {
		return nil
	}

	if err := h.Init(ctx, &msg); err != nil {
		return nil
	}

	h.Handle(ctx)
	return nil
}
