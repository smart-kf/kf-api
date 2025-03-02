package kffrontend

import (
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/endpoints/http/vo/kffrontend"
)

type QRCodeController struct {
	BaseController
}

func (c *QRCodeController) Scan(ctx *gin.Context) {
	var req kffrontend.QRCodeScanRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	return
}
