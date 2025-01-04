package kfbackend

import (
	"context"

	"github.com/make-money-fast/captcha"

	"github.com/smart-fm/kf-api/infrastructure/caches"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type LoginRequest struct {
	CardID      string `json:"cardID" binding:"required" doc:"卡密id" validate:"required"`
	Password    string `json:"password" doc:"密码可选"`
	CaptchaID   string `json:"captchaId" doc:"验证码id,通过验证码接口获取"`
	CaptchaCode string `json:"captchaCode" doc:"验证码"`
}

func (r *LoginRequest) Validate(ctx context.Context) error {
	if caches.BillSettingCacheInstance.IsVerifyCodeEnable() {
		if r.CaptchaID == "" && r.CaptchaCode == "" {
			return xerrors.NewParamsErrors("请输入验证码")
		}

		if !captcha.VerifyString(ctx, r.CaptchaID, r.CaptchaCode) {
			return xerrors.NewParamsErrors("验证码错误")
		}
	}
	return nil
}

type LoginResponse struct {
	Token  string `json:"token" doc:"用户的token"`
	Notice string `json:"notice,omitempty" doc:"公告通知"`
}
