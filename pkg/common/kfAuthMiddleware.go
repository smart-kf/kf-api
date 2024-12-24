package common

import (
	"context"
	"errors"
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

var kfCardIDKey = "kf-card-key"

func KFAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := request.ParseFromRequest(ctx.Request, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (any, error) {
			return []byte(config.GetConfig().JwtKey), nil
		})

		var bc BaseController
		if err != nil {
			bc.Error(ctx, xerrors.AuthError)
			return
		}

		if !token.Valid {
			bc.Error(ctx, xerrors.AuthError)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		cardID := claims["cardId"].(string)

		reqCtx := context.WithValue(ctx.Request.Context(), kfCardIDKey, cardID)
		ctx.Request = ctx.Request.WithContext(reqCtx)

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
	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JwtKey), nil
	})
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
