package kfbackend

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/captcha"

	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
)

type PublicController struct {
	BaseController
}

func (c *PublicController) GetCaptchaId(ctx *gin.Context) {
	id := captcha.NewLen(ctx.Request.Context(), 4)
	c.Success(
		ctx, kfbackend.GetQRCodeIDResponse{
			id,
		},
	)
}
