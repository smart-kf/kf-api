package kfbackend

import (
	"context"

	"github.com/smart-fm/kf-api/endpoints/common"
)

type LoginRequest struct {
	CardID      string `json:"cardID" binding:"required" doc:"卡密id" validate:"required"`
	Password    string `json:"password" doc:"密码可选"`
	CaptchaID   string `json:"captchaId" doc:"验证码id,通过验证码接口获取"`
	CaptchaCode string `json:"captchaCode" doc:"验证码"`
}

func (r *LoginRequest) Validate(ctx context.Context) error {
	// if caches.BillSettingCacheInstance.IsVerifyCodeEnable() {
	// 	if r.CaptchaID == "" && r.CaptchaCode == "" {
	// 		return xerrors.NewParamsErrors("请输入验证码")
	// 	}
	//
	// 	if !captcha.VerifyString(ctx, r.CaptchaID, r.CaptchaCode) {
	// 		return xerrors.NewParamsErrors("验证码错误")
	// 	}
	// }
	return nil
}

type LoginResponse struct {
	Token     string `json:"token" doc:"用户的token"`
	Notice    string `json:"notice,omitempty" doc:"公告通知"`
	CdnDomain string `json:"cdnDomain" doc:"静态资源域名"`
}

type LogRequest struct {
	common.PageRequest
	BeginTime int64  `json:"beginTime" form:"beginTime" doc:"开始时间 毫秒"`
	EndTime   int64  `json:"endTime" form:"endTime" doc:"结束时间 毫秒"`
	Function  string `json:"function" form:"function" doc:"操作类型枚举:客户,欢迎语,快速发送,智能回复,二维码,设置,话术复制"`
}

type LogResponse struct {
	Id         int64  `json:"id"`
	HandleFunc string `json:"function" gorm:"column:handle_func;type:varchar(255)" doc:"操作类型"` // 操作类型
	Content    string `json:"content" gorm:"column:content;longtext;" doc:"内容"`                // 操作内容
	Ip         string `json:"ip" gorm:"column:ip" doc:"ip"`
	CreateTime int64  `json:"createTime" gorm:"createTime" doc:"创建时间：毫秒"`
}

type ListLogResponse struct {
	List  []*LogResponse `json:"list"`
	Total int64          `json:"total"`
}
