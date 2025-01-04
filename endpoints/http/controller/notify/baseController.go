package notify

import (
	"github.com/smart-fm/kf-api/endpoints/http/middleware"
)

type BaseController struct {
	middleware.BaseController
}

func NewBaseController() *BaseController {
	return &BaseController{}
}
