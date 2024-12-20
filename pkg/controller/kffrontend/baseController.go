package kffrontend

import (
	"std-api/pkg/common"
)

type BaseController struct {
	common.BaseController
}

func NewBaseController() *BaseController {
	return &BaseController{}
}
