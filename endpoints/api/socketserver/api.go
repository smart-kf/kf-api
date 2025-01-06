package socketserver

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/smart-fm/kf-api/config"
)

var (
	SocketServerClient *socketServerClient
	o                  sync.Once
)

type socketServerClient struct {
	client *resty.Client
}

func NewSocketServerClient() *socketServerClient {
	o.Do(
		func() {
			r := resty.NewWithClient(
				&http.Client{
					Timeout: time.Duration(config.GetConfig().HttpClient.Timeout) * time.Second,
				},
			)
			SocketServerClient = &socketServerClient{
				r,
			}
		},
	)
	return SocketServerClient
}

type PushMessageRequest struct {
	SessionId string `json:"sessionId"`
	Event     string `json:"event"`
	Data      string `json:"data"`
}

func (r *PushMessageRequest) SetData(v interface{}) {
	data, _ := json.Marshal(v)
	r.Data = string(data)
}

func (c *socketServerClient) PushMessage(ctx context.Context, msg *PushMessageRequest) error {
	url := config.GetConfig().HttpClient.SocketServerClient + "/api/push"
	_, err := c.client.R().SetBody(msg).Post(url)
	if err != nil {
		return err
	}

	return nil
}
