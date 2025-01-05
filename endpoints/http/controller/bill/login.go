package bill

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/http/vo/bill"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

func (c *BaseController) Login(ctx *gin.Context) {
	var request bill.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.Error(ctx, err)
		return
	}

	if err := request.Validate(); err != nil {
		c.Error(ctx, err)
		return
	}

	reqCtx := ctx.Request.Context()

	var br repository.BillAccountRepository

	bac, ok, err := br.FindOneByUsername(reqCtx, request.Username)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	if !ok {
		c.Error(ctx, xerrors.NewCustomError("账号或密码错误"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(bac.Password), []byte(request.Password)); err != nil {
		c.Error(ctx, xerrors.NewCustomError("账号或密码错误"))
		return
	}

	// 登录正确. 设置jwt - token
	xlogger.Info(ctx, "用户登录系统", xlogger.Any("req", request))

	j := jwt.New(jwt.SigningMethodHS256)
	j.Claims = jwt.MapClaims{
		"id":       bac.ID,
		"username": bac.Username,
	}
	sign, err := j.SignedString([]byte(config.GetConfig().JwtKey))
	if err != nil {
		c.Error(ctx, err)
		return
	}

	c.Success(
		ctx, bill.LoginResponse{
			Token: sign,
		},
	)
}
