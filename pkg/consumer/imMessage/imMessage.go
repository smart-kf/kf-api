package imMessage

import (
	"context"
	"encoding/binary"

	"github.com/IBM/sarama"
	xlogger "github.com/clearcodecn/log"
)

type ImMessageConsumer struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	stopChan   chan struct{}
}

func NewImMessageConsumer(stopChan chan struct{}) (*ImMessageConsumer, error) {
	var c = &ImMessageConsumer{
		stopChan: stopChan,
	}
	return c, nil
}

func (c *ImMessageConsumer) Consume() error {
	// conf := config.GetConfig().Kafka
	// config := sarama.NewConfig()
	// config.Version = sarama.V1_0_0_0
	// config.Consumer.Return.Errors = true
	// config.Consumer.Offsets.Initial = sarama.OffsetNewest
	// config.Consumer.Group.Session.Timeout = 20 * time.Second
	// config.Consumer.Group.Heartbeat.Interval = 6 * time.Second
	// config.Consumer.MaxProcessingTime = 500 * time.Millisecond
	// config.Net.DialTimeout = time.Second * 10
	//
	// cg, err := sarama.NewConsumerGroup(conf.Addrs, conf.ImMessageGroup, config)
	// if err != nil {
	// 	return err
	// }
	// ctx, cancel := context.WithCancel(context.Background())
	//
	// c.ctx = ctx
	// c.cancelFunc = cancel
	// return cg.Consume(ctx, []string{conf.ImMessageTopic}, c)
	return nil
}

func (c *ImMessageConsumer) Setup(session sarama.ConsumerGroupSession) error {
	xlogger.Info(context.Background(), "ImMessageConsumer-Setup")
	return nil
}

func (c ImMessageConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	xlogger.Info(c.ctx, "ImMessageConsumer Cleanup")
	if c.cancelFunc != nil {
		c.cancelFunc()
	}
	return nil
}

func (c ImMessageConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			xlogger.Info(c.ctx, "Receive new Message:", xlogger.Any("val", string(msg.Value)))
			operation, msgBody := unPackMessage(msg.Value)
			if err := handleOperation(c.ctx, operation, msgBody); err != nil {
				xlogger.Error(c.ctx, "handleMessage failed", xlogger.Err(err))
				session.MarkMessage(msg, "error")
				continue
			}
			session.MarkMessage(msg, "")
		case <-c.ctx.Done():
			return nil
		case <-c.stopChan:
			return nil
		}
	}
}

// unPackMessage 解包消息
func unPackMessage(v []byte) (int32, []byte) {
	var op uint32
	op = binary.LittleEndian.Uint32(v[:4])

	return int32(op), v[4:]
}
