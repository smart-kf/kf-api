package bill

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"std-api/config"
	"std-api/pkg/repository"
	"std-api/pkg/xerrors"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required" doc:"用户名"`
	Password string `json:"password" binding:"required" doc:"密码"`
}

func (r *LoginRequest) Validate() error {
	if len(r.Username) == 0 || len(r.Password) == 0 {
		return xerrors.NewParamsErrors("账号或密码不能为空")
	}
	return nil
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (c *BillController) Login(ctx *gin.Context) {
	var request LoginRequest
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
	xlogger.Info("用户登录系统", xlogger.Any("req", request))

	j := jwt.New(jwt.SigningMethodHS256)
	j.Claims = jwt.MapClaims{
		"username": bac.Username,
	}
	sign, err := j.SignedString([]byte(config.GetConfig().JwtKey))
	if err != nil {
		c.Error(ctx, err)
		return
	}

	c.Success(ctx, LoginResponse{
		Token: sign,
	})
}
