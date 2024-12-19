package utils

import (
	"github.com/gin-gonic/gin"
	"net"
	"strings"
)

func ClientIP(ctx *gin.Context) string {
	if ip := ctx.Request.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	if ip := ctx.Request.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(ctx.Request.RemoteAddr))
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}
