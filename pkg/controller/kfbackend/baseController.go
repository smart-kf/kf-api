package kfbackend

import (
	"github.com/smart-fm/kf-api/pkg/common"
)

type BaseController struct {
	common.BaseController
}

func NewBaseController() *BaseController {
	return &BaseController{}
}
