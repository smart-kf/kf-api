package middleware

import (
	"context"
	"fmt"
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/infrastructure/redis"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

var kfCardIDKey = "kf-card-key"
var kfTokenKey = "kf-token-key"

func KFAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		// 空的也可以请求，由控制器生成 token 返回前端.
		if token == "" {
			ctx.Next()
			return
		}
		redisClient := redis.GetRedisClient()
		cardID, err := redisClient.Get(ctx.Request.Context(), fmt.Sprintf("kfbe.%s", token)).Result()
		var bc BaseController
		if err != nil {
			bc.Error(ctx, xerrors.AuthError)
			return
		}
		reqCtx := context.WithValue(ctx.Request.Context(), kfCardIDKey, cardID)
		reqCtx = context.WithValue(reqCtx, kfTokenKey, token)
		ctx.Request = ctx.Request.WithContext(reqCtx)
		// TODO:: 验证 card 的有效期, 这里给token续期.
		redisClient.Set(ctx.Request.Context(), fmt.Sprintf("kfbe.%s", token), cardID, 7*24*time.Hour)
		ctx.Next()
	}
}

func GetKFCardID(ctx *gin.Context) string {
	acc, ok := ctx.Request.Context().Value(kfCardIDKey).(string)
	if ok {
		return acc
	}
	xlogger.Warn(ctx.Request.Context(), "GetKFCardID-failed, 未从context中获取到cardID,请检查路由设置")
	return ""
}

func GetKFToken(ctx *gin.Context) string {
	acc, ok := ctx.Request.Context().Value(kfTokenKey).(string)
	if ok {
		return acc
	}
	xlogger.Warn(ctx.Request.Context(), "GetKFCardID-failed, 未从context中获取到kfToken,请检查路由设置")
	return ""
}

func VerifyKFBackendToken(ctx context.Context, token string) error {
	redisClient := redis.GetRedisClient()
	_, err := redisClient.Get(ctx, fmt.Sprintf("kfbe.%s", token)).Result()
	if err != nil {
		return err
	}
	return nil
}

// TODO:: VerifyKFToken 前台用户的token
func VerifyKFToken(ctx context.Context, s string) error {
	return nil
}
