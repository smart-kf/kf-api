package kfbackend

import (
	"unicode/utf8"

	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type UpsertWelcomeMsgRequest struct {
	Id      int    `json:"id" doc:"id: 创建传0或者不传"`
	Content string `json:"content" doc:"内容: text: 不超过255个字符, video: url, image: url" binding:"required"`
	Type    string `json:"type" doc:"text,video,image" binding:"required"`
	Sort    int    `json:"sort" doc:"排序编号,欢迎语、快捷回复需要"`
	Enable  bool   `json:"enable" doc:"是否启用"`
	Keyword string `json:"keyword" doc:"关键词,智能回复需要必填"`
	MsgType string `json:"msgType" binding:"required"  doc:"快捷回复=quick_reply, 欢迎语=welcome_msg, 智能回复=smart_msg
" validate:"oneof=quick_reply welcome_msg smart_msg"`
}

func (r *UpsertWelcomeMsgRequest) Validate() error {
	if utf8.RuneCountInString(r.Content) > 255 {
		return xerrors.NewParamsErrors("参数错误")
	}
	switch r.Type {
	case "text", "image", "video":
	default:
		return xerrors.NewParamsErrors("参数错误")
	}
	if r.MsgType == "smart_msg" && r.Keyword == "" {
		return xerrors.NewParamsErrors("关键词错误")
	}
	return nil
}

type DeleteWelcomeRequest struct {
	Id int `json:"id" doc:"主键id" binding:"required"`
}

type KfWelcomeMessageResp struct {
	Id      int64  `json:"id"`
	Content string `json:"content" gorm:"type:text"`
	Type    string `json:"type" gorm:"type:varchar(255)"`
	Sort    int    `json:"sort"`   // 排序
	Enable  bool   `json:"enable"` // 是否启用.
	Keyword string `json:"keyword"`
}

type KfWelcomeMessageListResp struct {
	List  []*KfWelcomeMessageResp `json:"list"`
	Total int64                   `json:"total"`
}

type ListAllRequest struct {
	common.PageRequest
	MsgType string `form:"msgType" json:"msgType" binding:"required"  doc:"快捷回复=quick_reply, 
欢迎语=welcome_msg 智能回复=smart_msg" validate:"oneof=quick_reply welcome_msg smart_msg"`
}
