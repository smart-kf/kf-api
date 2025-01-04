package consumer

import (
	"fmt"

	"github.com/nsqio/go-nsq"
)

// MessageConsumer 消息消费者.
type MessageConsumer struct{}

func (m *MessageConsumer) HandleMessage(message *nsq.Message) error {
	fmt.Println("receive a new message --->", string(message.Body))
	return nil
}
