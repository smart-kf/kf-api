package wsclient

import (
	"context"
	"encoding/json"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/dto"
	"github.com/smart-fm/kf-api/pkg/httpclient"
)

type PushMessageRequest struct {
	SessionId  string   `json:"sessionId"`
	SessionIds []string `json:"sessionIds"`
	Event      string   `json:"event"`
	Data       string   `json:"data"`
}

type WsClient struct{}

func (WsClient) Push(ctx context.Context, event string, message *dto.Message, sessionIds ...string) error {
	url := config.GetConfig().HttpClient.SocketServerClient + "/api/push"
	msg := PushMessageRequest{
		Event: event,
	}
	if len(sessionIds) == 1 {
		msg.SessionId = sessionIds[0]
	} else {
		msg.SessionIds = sessionIds
	}

	msgBody, _ := json.Marshal(message)
	msg.Data = string(msgBody)

	var out = make(map[string]interface{})

	err := httpclient.Post(ctx, url, msg, &out)
	if err != nil {
		return err
	}
	return nil
}
