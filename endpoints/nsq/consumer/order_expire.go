package consumer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nsqio/go-nsq"

	"github.com/smart-fm/kf-api/domain/service/orders"
)

// MessageConsumer 消息消费者.
type OrderExpireConsumer struct{}

func (m *OrderExpireConsumer) HandleMessage(message *nsq.Message) error {
	fmt.Println("receive a new message --->", string(message.Body))
	// 创建消息，并且给客户端回复消息id.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_ = ctx

	var orderNo string = string(message.Body)

	if strings.HasPrefix(orderNo, "X") {
		orders.ExpireDomainOrder(ctx, orderNo)
	}
	if strings.HasPrefix(orderNo, "N") {
		orders.CardOrderExpire(ctx, orderNo)
	}
	return nil
}
