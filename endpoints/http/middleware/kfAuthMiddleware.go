package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/factory"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

func KFBeAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var bc BaseController
		token := ctx.GetHeader("Authorization")
		// 空的也可以请求，由控制器生成 token 返回前端.
		if token == "" {
			bc.Error(ctx, xerrors.AuthError)
			return
		}
		cardId, err := caches.KfAuthCacheInstance.GetBackendToken(ctx.Request.Context(), token)
		if err != nil {
			bc.Error(ctx, xerrors.AuthError)
			return
		}
		newCtx := common.WithKfCardID(ctx.Request.Context(), cardId)
		newCtx = common.WithKFToken(newCtx, token)
		ctx.Request = ctx.Request.WithContext(newCtx)
		// TODO:: 验证 card 的有效期, 这里给token续期.
		ctx.Next()
	}
}

// KFFeAuthMiddleware 如果redis中token过期，需要去数据库查询，然后设置到redis中.
func KFFeAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var bc BaseController
		token := ctx.GetHeader("Authorization")
		// 空的也可以请求，由控制器生成 token 返回前端.
		if token == "" {
			ctx.Next()
			return
		}
		// 可能造成数据库击穿.
		cardId, err := caches.KfAuthCacheInstance.GetFrontToken(ctx.Request.Context(), token)
		if err != nil {
			if strings.Contains(err.Error(), "token not found") {
				cardMainId, err := factory.FactoryParseUserToken(token)
				if err != nil {
					bc.Error(ctx, xerrors.AuthError)
					return
				}
				cardId, err = caches.KfCardCacheInstance.GetCardIDByMainID(ctx.Request.Context(), cardMainId)
				if err != nil {
					bc.Error(ctx, xerrors.AuthError)
					return
				}
				_, err = caches.KfUserCacheInstance.GetDBUser(ctx, cardId, token)
				if err != nil {
					bc.Error(ctx, xerrors.AuthError)
					return
				}
				caches.KfAuthCacheInstance.SetFrontToken(ctx.Request.Context(), token, cardId)
			} else {
				bc.Error(ctx, xerrors.AuthError)
				return
			}
		}
		newCtx := common.WithKfCardID(ctx.Request.Context(), cardId)
		newCtx = common.WithKFToken(newCtx, token)
		ctx.Request = ctx.Request.WithContext(newCtx)
		// TODO:: 验证 card 的有效期, 这里给token续期.
		ctx.Next()
	}
}

// VerifyKFToken 前台用户的token
func VerifyKFToken(ctx context.Context, token string) error {
	_, err := caches.KfAuthCacheInstance.GetFrontToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}
