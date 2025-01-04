package dev

import (
	"github.com/clearcodecn/swaggos"

	"github.com/smart-fm/kf-api/endpoints/common"
)

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/dev").Tag("开发使用的接口")
	bg.Post("/push").Body(PushMsgRequest{}).JSON(common.EmptyResponse{}).Description("模拟推送消息")
}
