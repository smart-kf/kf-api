package controller

import (
	"github.com/gin-gonic/gin"
	"std-api/config"
	"std-api/pkg/xerrors"
)

type BaseResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	DebugInfo any    `json:"debug_info,omitempty"`
	Data      any    `json:"data"`
}

type BaseController struct{}

func (b *BaseController) Success(ctx *gin.Context, obj any) {
	ctx.JSON(200, &BaseResponse{
		Code: 200,
		Data: obj,
	})
}

func (b *BaseController) Error(ctx *gin.Context, err error) {
	var rsp BaseResponse
	if myError, ok := xerrors.IsError(err); ok {
		rsp.Code = myError.Code
		rsp.Message = myError.Msg

		ctx.JSON(200, rsp)
		return
	}

	rsp.Code = 500
	rsp.Message = "internal server error"
	if config.GetConfig().Debug {
		rsp.DebugInfo = err
	}
	ctx.JSON(200, rsp)
}
