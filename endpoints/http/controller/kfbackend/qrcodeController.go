package kfbackend

import (
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
)

type QRCodeController struct {
	BaseController
}

func (c *QRCodeController) List(ctx *gin.Context) {
	var req kfbackend.QRCodeRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
}
