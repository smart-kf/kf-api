package common

import (
	"context"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
)

// Context, kf 前台、kf后台通用，存储 cardId 信息, 做了 *gin.Context 兼容

var kfCardIDKey = "kf-card-key"
var kfTokenKey = "kf-token-key"

func WithKfCardID(ctx context.Context, cardId string) context.Context {
	return context.WithValue(ctx, kfCardIDKey, cardId)
}

func GetKFCardID(ctx context.Context) string {
	ginCtx, ok := ctx.(*gin.Context)
	if ok {
		ctx = ginCtx.Request.Context()
	}

	acc, ok := ctx.Value(kfCardIDKey).(string)
	if ok {
		return acc
	}
	xlogger.Warn(ctx, "GetKFCardID-failed, 未从context中获取到cardID,请检查路由设置")
	return ""
}

func WithKFToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, kfTokenKey, token)
}

func GetKFToken(ctx context.Context) string {
	ginCtx, ok := ctx.(*gin.Context)
	if ok {
		ctx = ginCtx.Request.Context()
	}
	acc, ok := ctx.Value(kfTokenKey).(string)
	if ok {
		return acc
	}
	xlogger.Warn(ctx, "GetKFCardID-failed, 未从context中获取到kfToken,请检查路由设置")
	return ""
}
