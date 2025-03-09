package dto

import (
	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type Message struct {
	Event       string `json:"event"`
	Platform    string `json:"platform,omitempty"`  // platform  来源平台.
	SessionId   string `json:"sessionId,omitempty"` // sessionId
	Token       string `json:"token,omitempty"`     // token: 前端/后台的token，一定有.
	MsgType     string `json:"msgType"`             // text || image || video
	MsgId       string `json:"msgId"`               // 消息id
	GuestName   string `json:"guestName"`           // 客户名称
	GuestAvatar string `json:"guestAvatar"`         // 客户头像
	GuestId     string `json:"guestId"`             // 客户id: 后台推送前台才有.
	Content     string `json:"content"`             // 具体消息内容
	KfId        string `json:"kfId"`                // 客服id
	IsKf        int    `json:"isKf"`                // 1=客服，2=粉丝.
}

func (m Message) FromBackend() bool {
	return m.Platform == constant.PlatformKfBe
}

func NewGuestOnlineMessage(guestName string, guestId string, avatar string, cardId string, sessionId string) *Message {
	return &Message{
		Event:       constant.EventOnline,
		Platform:    constant.PlatformKfBe,
		SessionId:   sessionId,
		GuestName:   guestName,
		GuestAvatar: avatar,
		GuestId:     guestId,
		KfId:        cardId,
	}
}

func NewGuestOfflineMessage(guestId string, cardId string, sessionId string) *Message {
	return &Message{
		Event:     constant.EventOnline,
		Platform:  constant.PlatformKfBe,
		SessionId: sessionId,
		GuestId:   guestId,
		KfId:      cardId,
	}
}

func NewMessage(oldMessage *Message, toPlatform string) *Message {
	return &Message{
		Event:       constant.EventMessage,
		Platform:    toPlatform,
		MsgType:     oldMessage.MsgType,
		MsgId:       oldMessage.MsgId,
		GuestName:   oldMessage.GuestName,
		GuestAvatar: config.GetConfig().Web.CdnHost + oldMessage.GuestAvatar,
		GuestId:     oldMessage.Token,
		Content:     oldMessage.Content,
		KfId:        oldMessage.KfId,
		IsKf:        oldMessage.IsKf,
	}
}

func NewPushMessage(msgType string, msgId string, content string, user *dao.KfUser) *Message {
	return &Message{
		Event:       constant.EventMessage,
		Platform:    constant.PlatformKfFe,
		MsgType:     msgType,
		MsgId:       msgId,
		GuestName:   user.NickName,
		GuestAvatar: config.GetConfig().Web.CdnHost + user.Avatar,
		GuestId:     user.UUID,
		Content:     content,
		IsKf:        constant.IsKf,
	}
}

//
// type Online struct {
// 	Event       string `json:"event"`
// 	Platform    string `json:"platform,omitempty"`  // platform  来源平台.
// 	SessionId   string `json:"sessionId,omitempty"` // sessionId
// 	Token       string `json:"token,omitempty"`     // token
// 	GuestName   string `json:"guestName"`           // 客户名称
// 	GuestAvatar string `json:"guestNvatar"`         // 客户头像
// 	GuestId     string `json:"guestId"`             // 客户id
// 	IsKf        int    `json:"isKf"`                // 1=客服，2=粉丝.
// 	KfId        string `json:"kfId"`                // 客服id
// }
