package dto

import "github.com/smart-fm/kf-api/endpoints/common/constant"

type MessageBase struct {
	Event     string `json:"event"`
	Platform  string `json:"platform,omitempty"`  // platform  来源平台.
	SessionId string `json:"sessionId,omitempty"` // sessionId
	Token     string `json:"token,omitempty"`     // token
}

type Message struct {
	MessageBase
	MsgType     string `json:"msgType"`     // text || image || video
	MsgId       string `json:"msgId"`       // 消息id
	GuestName   string `json:"guestName"`   // 客户名称
	GuestAvatar string `json:"guestAvatar"` // 客户头像
	GuestId     string `json:"guestId"`     // 客户id
	Content     string `json:"content"`     // 具体消息内容
	KfId        string `json:"kfId"`        // 客服id
	IsKf        int    `json:"isKf"`        // 1=客服，2=粉丝.
}

func (m MessageBase) FromBackend() bool {
	return m.Platform == constant.PlatformKfBe
}

type Online struct {
	MessageBase
	GuestName   string `json:"guestName"`   // 客户名称
	GuestAvatar string `json:"guestNvatar"` // 客户头像
	GuestId     string `json:"guestId"`     // 客户id
	IsKf        int    `json:"isKf"`        // 1=客服，2=粉丝.
	KfId        string `json:"kfId"`        // 客服id
}
