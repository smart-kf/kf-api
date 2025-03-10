package sdk

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type Address struct {
	Id      int    `json:"id" gorm:"primaryKey"`
	AppId   string `json:"app_id"`
	Address string `json:"address"`
	Enable  bool   `json:"enable"`
	Remark  string `json:"remark"`
}

type ListAddressRequest struct {
	AppId string `json:"appId"`
}

func (ListAddressRequest) Url() string {
	return "/api/v1/admin/address/list"
}

type UpsertAddressRequest struct {
	Id      int    `json:"id"`
	AppId   string `json:"appId"`
	Enable  bool   `json:"enable"`
	Address string `json:"address"`
	Remark  string `json:"remark"`
}

func (UpsertAddressRequest) Url() string {
	return "/api/v1/admin/address/upsert"
}

type DeleteAddressRequest struct {
	Id    int    `json:"id"`
	AppId string `json:"appId"`
}

func (DeleteAddressRequest) Url() string {
	return "/api/v1/admin/address/del"
}

type UpdateTronRequest struct {
	ApiKey      string `json:"apiKey" binding:"required"` // 可能有多个，以 , 分割.
	Proxy       string `json:"proxy"`
	Timeout     int    `json:"timeout"`
	TronNetwork string `json:"tron_network" binding:"required"`
	CronSecond  int    `json:"cron_second" binding:"required"` // 定时任务执行秒数间隔
}

func (UpdateTronRequest) Url() string {
	return "/api/v1/admin/tron/update"
}

type GetTronRequest struct {
}

func (GetTronRequest) Url() string {
	return "/api/v1/admin/tron/get"
}

type ManagementAdapter interface {
	Url() string
}

func (c UsdtPaymentClient) Management(ctx context.Context, req ManagementAdapter) ([]byte, error) {
	r := resty.New().R().SetContext(ctx).SetBody(req).SetHeader("Authorization", c.token)
	res, err := r.Post(c.host + req.Url())
	if err != nil {
		return nil, err
	}
	return res.Body(), nil
}
