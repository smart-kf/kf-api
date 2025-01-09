package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"

	"github.com/smart-fm/kf-api/infrastructure/httpClient/socketserver"
)

// MessageConsumer 消息消费者.
type MessageConsumer struct{}

func (m *MessageConsumer) HandleMessage(message *nsq.Message) error {
	fmt.Println("receive a new message --->", string(message.Body))
	// 创建消息，并且给客户端回复消息id.

	var msg Message
	err := json.Unmarshal(message.Body, &msg)
	if err != nil {
		return err
	}

	client := socketserver.NewSocketServerClient()
	switch msg.Event {
	case EventSessionId:
		// 向客户端推送一个sessionId 消息.
		// TODO:: 更好的实现方式
		var req = socketserver.PushMessageRequest{
			SessionId: msg.SessionId,
			Event:     EventSessionId,
		}
		req.SetData(
			map[string]string{
				"sessionId": msg.SessionId,
			},
		)
		client.PushMessage(context.Background(), &req)
	case EventMessage:
		var req = socketserver.PushMessageRequest{
			SessionId: msg.SessionId,
			Event:     EventMessageAck,
		}
		req.SetData(
			&Message{
				Event:     EventMessageAck,
				MsgId:     uuid.New().String(),
				Platform:  msg.Platform,
				SessionId: msg.SessionId,
			},
		)
		client.PushMessage(context.Background(), &req)
	}

	return nil
}

const (
	PlatformKF        = "kf"         // 前台
	PlatformKFBackend = "kf-backend" // 客服后台
)

const (
	EventSessionId  = "sessionId"  // 新建连接、初始化sessionId事件
	EventDisConnect = "disConnect" // 连接断开事件
	EventMessage    = "message"    // 发送消息事件
	EventMessageAck = "messageAck" // 发送消息事件
	EventOnline     = "online"     // 上线事件
	EventOffline    = "offline"    // 下线事件
)

type Message struct {
	Event     string `json:"event"`
	Platform  string `json:"platform,omitempty"`  // platform
	SessionId string `json:"sessionId,omitempty"` // sessionId
	Token     string `json:"token,omitempty"`     // token

	MsgType     string `json:"msgType"`     // text || image || video
	MsgId       string `json:"msgId"`       // 消息id
	GuestName   string `json:"guestName"`   // 客户名称
	GuestAvatar string `json:"guestNvatar"` // 客户头像
	GuestId     string `json:"guestId"`     // 客户id
	Content     string `json:"content"`     // 具体消息内容
	KfId        string `json:"kfId"`        // 客服id
	IsKf        int    `json:"isKf"`        // 1=客服，2=粉丝.
}
