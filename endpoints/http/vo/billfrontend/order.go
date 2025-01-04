package billfrontend

import (
	"context"

	"github.com/make-money-fast/captcha"

	"github.com/smart-fm/kf-api/pkg/caches"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type CreateOrderRequest struct {
	CaptchaId   string `json:"captchaId" doc:"验证码id"`
	CaptchaCode string `json:"captchaCode" doc:"验证码内容"`
}

func (r *CreateOrderRequest) Validate(ctx context.Context) error {
	if caches.BillSettingCacheInstance.IsVerifyCodeEnable() {
		if !captcha.VerifyString(ctx, r.CaptchaId, r.CaptchaCode) {
			return xerrors.NewParamsErrors("验证码错误")
		}
	}
	return nil
}

type CreateOrderResponse struct {
	OrderNo string `json:"orderNo" doc:"订单号"`
}
