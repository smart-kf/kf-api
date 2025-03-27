package kflog

import "github.com/smart-fm/kf-api/endpoints/common/constant"

var Functions = map[string]string{
	constant.WelcomeMsg: "欢迎语",
	constant.QuickReply: "快速发送",
	constant.SmartMsg:   "智能回复",
	"login":             "登录",
}

func MaskContent(s string) string {
	if len(s) == 4 {
		return "****"
	}
	return s[:4] + "****"
}
