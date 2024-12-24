package common

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/pkg/db"
	"github.com/smart-fm/kf-api/pkg/xerrors"
	"gorm.io/gorm"
)

var billAccountInfoKey = "ctx-key-bill-account"

func BillAuthMiddleware() gin.HandlerFunc {
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

		username, _ := claims["username"].(string)
		uid, _ := claims["id"].(float64)

		reqCtx := context.WithValue(ctx.Request.Context(), billAccountInfoKey, &db.BillAccount{
			Model:    gorm.Model{ID: uint(uid)},
			Username: username,
		})
		ctx.Request = ctx.Request.WithContext(reqCtx)

		ctx.Next()
	}
}

func GetBillAccount(ctx *gin.Context) db.BillAccount {
	acc, ok := ctx.Request.Context().Value(billAccountInfoKey).(*db.BillAccount)
	if ok {
		return *acc
	}
	return db.BillAccount{Username: "unknown"}
}

func VerifyKFBackendToken(s string) (*db.BillAccount, error) {
	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token invalid")
	}

	claims := token.Claims.(jwt.MapClaims)

	username, _ := claims["username"].(string)
	uid, _ := claims["id"].(float64)
	return &db.BillAccount{
		Model:    gorm.Model{ID: uint(uid)},
		Username: username,
	}, err
}
