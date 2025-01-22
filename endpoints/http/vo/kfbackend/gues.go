package kfbackend

import "errors"

type GetKfUserInfoRequest struct {
	UUID string `json:"uuid" doc:"用户uuid" binding:"required" form:"uuid" query:"uuid"`
}

type GetKfUserInfoResponse struct {
	User
}

type UpdateUserInfoRequest struct {
	UUID       string `json:"uuid" binding:"required" doc:"用户uuid"`
	UpdateType string `json:"updateType" binding:"required" doc:"更新类型: userinfo=用户信息, block=拉黑, top=置顶"`
	Block      int    `json:"block" doc:"拉黑=1，取消拉黑=2"`
	Top        int    `json:"top" doc:"置顶=1，取消置顶=2"`
	RemarkName string `json:"remarkName" gorm:"column:remark_name" doc:"备注名称,尅可以为空"` // 备注名称.
	Mobile     string `json:"mobile" gorm:"column:mobile" doc:"手机号,可以为空"`            // 手机号
	Comments   string `json:"comments" gorm:"column:comments" doc:"备注信息,可以为空"`       // 备注信息
}

func (r *UpdateUserInfoRequest) Validate() error {
	switch r.UpdateType {
	case "userinfo":
	case "block":
	case "top":
	default:
		return errors.New("更新类型错误")
	}
	return nil
}

func (r *UpdateUserInfoRequest) IsUserInfo() bool {
	return r.UpdateType == "userinfo"
}

func (r *UpdateUserInfoRequest) IsBlock() bool {
	return r.UpdateType == "block"
}

func (r *UpdateUserInfoRequest) IsTop() bool {
	return r.UpdateType == "top"
}
