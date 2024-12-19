package bill

import "std-api/pkg/controller"

type BillController struct {
	controller.BaseController
}

func NewKfBaseController() *BillController {
	return &BillController{}
}
