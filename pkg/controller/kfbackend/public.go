package kfbackend

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/captcha"
)

type PublicController struct {
	BaseController
}

type GetQRCodeIDRequest struct{}

type GetQRCodeIDResponse struct {
	Id string `json:"id"`
}

func (c *PublicController) GetCaptchaId(ctx *gin.Context) {
	id := captcha.NewLen(ctx.Request.Context(), 4)
	c.Success(ctx, GetQRCodeIDResponse{
		id,
	})
}
