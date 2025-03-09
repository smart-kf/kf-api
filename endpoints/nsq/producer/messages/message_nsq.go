package messages

import (
	"context"
	"encoding/json"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/endpoints/nsq/producer"
)

func PushMessages(ctx context.Context, messages ...*dto.Message) error {
	for _, msg := range messages {
		body, _ := json.Marshal(msg)
		if err := producer.NSQProducer.Publish(config.GetConfig().NSQ.MessageTopic, body); err != nil {
			return err
		}
	}
	return nil
}
