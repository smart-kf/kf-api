package kfbackend

import "github.com/gin-gonic/gin"

type QRCodeController struct {
	BaseController
}

func (c *QRCodeController) List(ctx *gin.Context) {
	var req QRCodeRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
}
