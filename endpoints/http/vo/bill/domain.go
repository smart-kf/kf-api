package bill

import (
	"github.com/smart-fm/kf-api/endpoints/common"
)

type AddDomainRequest struct {
	TopName  string `json:"topName" binding:"required" validate:"required" doc:"顶级域名"`
	Status   int    `json:"status" binding:"required" validate:"required,oneof=1 2" doc:"域名状态: 1=正常，2=禁用"`
	IsPublic *bool  `json:"isPublic" binding:"required" validate:"required" doc:"true=公共域名,false=私有域名"`
}

func (r *AddDomainRequest) Validate() error {
	return nil
}

type ListDomainRequest struct {
	common.PageRequest
}

type ListDomainResponse struct {
	List  []*BillDomainResponse `json:"list" doc:"列表数据"`
	Total int64                 `json:"total"`
}

type BillDomainResponse struct {
	ID         int64  `json:"id"`
	TopName    string `json:"topName" doc:"域名"`
	IsPublic   bool   `json:"isPublic" doc:"是否公共域名"`
	IsBind     bool   `json:"isBind" doc:"是否有卡密绑定"`
	Status     int    `json:"status"  doc:"域名状态: 1=正常，2=禁用"`
	CreateTime int64  `json:"createTime" doc:"创建时间"`
}

type DeleteDomainRequest struct {
	ID int64 `json:"id" doc:"主键id" binding:"required" validate:"required"`
}
