package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type BaseResponse struct {
	Code      int    `json:"code" doc:"响应码:200=成功，400-499是业务错误, 500以上是服务器错误"`
	Message   string `json:"message" doc:"错误提示"`
	DebugInfo any    `json:"debug_info,omitempty" doc:"debug信息，开发环境才有，前端不要用这个字段做业务"`
	Data      any    `json:"data" doc:"数据信息"`
}

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	zhon := zhongwen.New()
	uni := ut.New(zhon, zhon)
	trans, _ = uni.GetTranslator("zh")
	validate = validator.New()
	err := zh.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		panic(err)
	}
}

type BaseController struct{}

func (b *BaseController) Success(ctx *gin.Context, obj any) {
	if obj == nil {
		obj = gin.H{}
	}
	ctx.JSON(
		200, &BaseResponse{
			Code: 200,
			Data: obj,
		},
	)
}

func (b *BaseController) Error(ctx *gin.Context, err error) {
	var rsp BaseResponse
	if myError, ok := xerrors.IsError(err); ok {
		rsp.Code = myError.Code
		rsp.Message = myError.Msg
		ctx.JSON(200, rsp)
		ctx.Abort()
		return
	}

	rsp.Code = 500
	rsp.Message = "internal server error"
	if config.GetConfig().Debug {
		rsp.DebugInfo = err.Error()
	}
	ctx.JSON(200, rsp)
	ctx.Abort()
}

func (b *BaseController) BindAndValidate(ctx *gin.Context, req any) bool {
	if err := ctx.ShouldBind(req); err != nil {
		b.Error(ctx, err)
		return false
	}
	err := validate.Struct(req)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			translations := validationErrors.Translate(trans)
			b.Error(ctx, xerrors.NewParamsValidateError(translations))
		}
		return false
	}
	if v, ok := req.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			b.Error(ctx, err)
			return false
		}
	}
	if v, ok := req.(interface {
		Validate(ctx context.Context) error
	}); ok {
		if err := v.Validate(ctx.Request.Context()); err != nil {
			b.Error(ctx, err)
			return false
		}
	}
	return true
}
