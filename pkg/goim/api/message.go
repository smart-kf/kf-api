package api

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type Message struct {
	Type string           `json:"type"`
	Data *json.RawMessage `json:"data"`
}

func FromBytes(v []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(v, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, err
}

type MessageTypeMsg struct {
	MsgType     string `json:"msgType"`     // 消息类型：text=文本，video=视频,image=图片
	MsgId       int64  `json:"msgId"`       // 消息id
	GuestId     int64  `json:"guestId"`     // 客户id
	GuestName   string `json:"guestName"`   // 客户名称
	GuestAvatar string `json:"guestAvatar"` // 客户头像
	KfName      string `json:"kfName"`      // 客服名称.
	KfAvatar    string `json:"kfAvatar"`    // 客服头像
	MsgTime     int64  `json:"msgTime"`     // 消息时间
	KfId        int64  `json:"kfId"`        // 客服id
	Content     string `json:"content"`     // 内容.
	IsKF        int    `json:"isKf"`        // 是否是客服：1=是客服消息，2=是客户消息.
}

func IsKFMessage(i int) bool {
	return i == 1
}

type MessageTypeRead struct {
	GuestId int64 `json:"guestId"` // 客户id
	KfId    int64 `json:"kfId"`    // 客服id
	IsKF    int   `json:"isKf"`    // 是否是客户
}

type MessageTypeOnline struct {
	GuestId int64 `json:"guestId"` // 客户id
	KfId    int64 `json:"kfId"`    // 客服id
	IsKF    int   `json:"isKf"`    // 是否是客户
}

type MessageTypeOffline struct {
	GuestId int64 `json:"guestId"` // 客户id
	KfId    int64 `json:"kfId"`    // 客服id
	IsKF    int   `json:"isKf"`    // 是否是客户
}

const (
	MessageTypeDefineMsg     = "msg"
	MessageTypeDefineRead    = "read"
	MessageTypeDefineOnline  = "online"
	MessageTypeDefineOffline = "offline"
)

func (m *Message) GetMsg() (*MessageTypeMsg, error) {
	if m.Type != MessageTypeDefineMsg {
		return nil, fmt.Errorf("invalid msgType: except: msg, get: %s", m.Type)
	}
	var data MessageTypeMsg
	err := json.Unmarshal(*m.Data, &data)
	if err != nil {
		return nil, errors.Wrap(err, "parse message failed"+m.Type)
	}
	return &data, nil
}

func (m *Message) GetRead() (*MessageTypeRead, error) {
	if m.Type != MessageTypeDefineRead {
		return nil, fmt.Errorf("invalid msgType: except: read, get: %s", m.Type)
	}
	var data MessageTypeRead
	err := json.Unmarshal(*m.Data, &data)
	if err != nil {
		return nil, errors.Wrap(err, "parse message failed"+m.Type)
	}
	return &data, nil
}

func (m *Message) GetOnline() (*MessageTypeOnline, error) {
	if m.Type != MessageTypeDefineOnline {
		return nil, fmt.Errorf("invalid msgType: except: online, get: %s", m.Type)
	}
	var data MessageTypeOnline
	err := json.Unmarshal(*m.Data, &data)
	if err != nil {
		return nil, errors.Wrap(err, "parse message failed"+m.Type)
	}
	return &data, nil
}

func (m *Message) GetOffline() (*MessageTypeOffline, error) {
	if m.Type != MessageTypeDefineOffline {
		return nil, fmt.Errorf("invalid msgType: except: offline, get: %s", m.Type)
	}
	var data MessageTypeOffline
	err := json.Unmarshal(*m.Data, &data)
	if err != nil {
		return nil, errors.Wrap(err, "parse message failed: "+m.Type)
	}
	return &data, nil
}
