package controller

import (
	"github.com/gin-gonic/gin"
	"std-api/config"
	"std-api/pkg/xerrors"
)

type BaseResponse struct {
	Code      int    `json:"code" doc:"响应码:200=成功，400-499是业务错误, 500以上是服务器错误"`
	Message   string `json:"message" doc:"错误提示"`
	DebugInfo any    `json:"debug_info,omitempty" doc:"debug信息，开发环境才有，前端不要用这个字段做业务"`
	Data      any    `json:"data" doc:"数据信息"`
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
	if config.GetConfig().Debug {
		rsp.DebugInfo = err
	}
	if myError, ok := xerrors.IsError(err); ok {
		rsp.Code = myError.Code
		rsp.Message = myError.Msg
		ctx.JSON(200, rsp)
		return
	}

	rsp.Code = 500
	rsp.Message = "internal server error"
	ctx.JSON(200, rsp)
}
