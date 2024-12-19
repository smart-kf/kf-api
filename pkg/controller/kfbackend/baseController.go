package kfbackend

import "std-api/pkg/controller"

type BaseController struct {
	controller.BaseController
}

func NewBaseController() *BaseController {
	return &BaseController{}
}
