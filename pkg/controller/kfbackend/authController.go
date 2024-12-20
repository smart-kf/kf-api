package kfbackend

import (
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	CardID   string `json:"cardID" binding:"required" doc:"卡密id" validate:"required"`
	Password string `json:"password,omitempty" doc:"密码可选" validate:"required"`
}

type LoginResponse struct {
	Notice string `json:"notice,omitempty" doc:"公告通知"`
}

type AuthController struct {
	BaseController
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

}
