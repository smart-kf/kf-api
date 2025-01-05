package billfrontend

import (
	"context"
)

type CreateOrderRequest struct {
	// CaptchaId   string `json:"captchaId" doc:"验证码id"`
	// CaptchaCode string `json:"captchaCode" doc:"验证码内容"`
	PackageId string `json:"packageId" doc:"套餐id" binding:"required" validate:"required"`
}

func (r *CreateOrderRequest) Validate(ctx context.Context) error {
	// if caches.BillSettingCacheInstance.IsVerifyCodeEnable() {
	// 	if !captcha.VerifyString(ctx, r.CaptchaId, r.CaptchaCode) {
	// 		return xerrors.NewParamsErrors("验证码错误")
	// 	}
	// }
	return nil
}

type CreateOrderResponse struct {
	OrderNo string `json:"orderNo" doc:"订单号"`
}
