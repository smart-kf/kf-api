package middleware

import (
	"context"
	"errors"
	"fmt"
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/infrastructure/redis"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

var kfCardIDKey = "kf-card-key"

func KFAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		redisClient := redis.GetRedisClient()
		cardID, err := redisClient.Get(ctx.Request.Context(), fmt.Sprintf("kfbe.%s", token)).Result()
		var bc BaseController
		if err != nil {
			bc.Error(ctx, xerrors.AuthError)
			return
		}
		reqCtx := context.WithValue(ctx.Request.Context(), kfCardIDKey, cardID)
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

func VerifyKFToken(s string) (string, error) {
	token, err := jwt.Parse(
		s, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().JwtKey), nil
		},
	)
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("token invalid")
	}

	claims := token.Claims.(jwt.MapClaims)
	cardID := claims["cardId"].(string)
	return cardID, err
}
