package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type UsdtPaymentClient struct {
	host    string
	token   string
	timeout time.Duration
}

func NewUsdtPaymentClient(host, token string, timeout time.Duration) *UsdtPaymentClient {
	c := &UsdtPaymentClient{
		host:    host,
		token:   token,
		timeout: timeout,
	}
	return c
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func NewParamsError(msg string) error {
	return &Error{
		Code:    400,
		Message: msg,
	}
}

func NewCodeError(code int, msg string) error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

type CreateOrderRequest struct {
	AppId       string `json:"app_id"`
	OrderId     string `json:"order_id"`
	Name        string `json:"name"`
	Amount      int64  `json:"amount"`
	FromAddress string `json:"from_address"`
	Expire      int    `json:"expire"` // 过期秒数
}

func (r *CreateOrderRequest) Validate() error {
	if r.AppId == "" {
		return NewParamsError("appid 不能为空")
	}
	if r.OrderId == "" {
		return NewParamsError("orderid 不能为空")
	}
	if r.Name == "" {
		return NewParamsError("name 不能为空")
	}
	if r.FromAddress == "" {
		return NewParamsError("fromAddress不能为空")
	}
	if r.Expire == 0 {
		return NewParamsError("过期时间不能为空")
	}
	return nil
}

type CreateOrderResponse struct {
	Data CreateOrderData `json:"data"`
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
}

type CreateOrderData struct {
	TradeId string `json:"trade_id"`
	PayUrl  string `json:"pay_url"`
}

func (c *UsdtPaymentClient) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderData, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	r := resty.New().R().SetContext(ctx).SetBody(req).SetHeader("Authorization", c.token)
	rsp, err := r.Post(fmt.Sprintf("%s/api/v1/order/create", c.host))
	if err != nil {
		return nil, NewCodeError(500, err.Error())
	}
	var cor CreateOrderResponse
	err = json.Unmarshal(rsp.Body(), &cor)
	if err != nil {
		return nil, err
	}
	if cor.Code != 0 {
		return nil, NewCodeError(cor.Code, cor.Msg)
	}
	return &cor.Data, nil
}

type TradeOrders struct {
	Id          int64      `gorm:"primary_key;AUTO_INCREMENT;comment:id"`
	OrderId     string     `gorm:"type:varchar(255);not null;unique;color:blue;comment:客户订单ID"`
	AppId       string     `gorm:"column:app_id;type:varchar(255)"`
	TradeId     string     `gorm:"type:varchar(255);not null;unique;color:blue;comment:本地订单ID"`
	TradeHash   string     `gorm:"type:varchar(64);default:'';unique;comment:交易哈希"`
	Amount      int64      `gorm:"default:0;comment:USDT交易数额,实际金额"`
	Money       int64      `gorm:"default:0;comment:订单交易金额,乘以1e6"`
	Address     string     `gorm:"type:varchar(34);not null;comment:收款地址"`
	FromAddress string     `gorm:"type:varchar(34);not null;default:'';comment:支付地址"`
	Status      int        `gorm:"type:tinyint(1);not null;default:0;comment:交易状态 1：等待支付 2：支付成功 3：订单过期"`
	Name        string     `gorm:"type:varchar(64);not null;default:'';comment:商品名称"`
	NotifyNum   int        `gorm:"type:int(11);not null;default:0;comment:回调次数"`
	NotifyState int        `gorm:"type:tinyint(1);not null;default:0;comment:回调状态 1：成功 0：失败"`
	ExpiredAt   time.Time  `gorm:"type:timestamp;not null;comment:订单失效时间"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;type:timestamp;not null;comment:创建时间"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime;type:timestamp;not null;comment:更新时间"`
	ConfirmedAt *time.Time `gorm:"type:timestamp;null;comment:交易确认时间"`
}

type QueryOrderRequest struct {
	AppId   string `json:"app_id"`
	OrderId string `json:"order_id"`
}

type QueryOrderResponse struct {
	Data TradeOrders `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

func (c *UsdtPaymentClient) QueryOrder(ctx context.Context, req *QueryOrderRequest) (*TradeOrders, error) {
	r := resty.New().R().SetContext(ctx).SetBody(req).SetHeader("Authorization", c.token)
	rsp, err := r.Post(fmt.Sprintf("%s/api/v1/order/query", c.host))
	if err != nil {
		return nil, NewCodeError(500, err.Error())
	}
	var cor QueryOrderResponse
	err = json.Unmarshal(rsp.Body(), &cor)
	if err != nil {
		return nil, err
	}
	if cor.Code != 0 {
		return nil, NewCodeError(cor.Code, cor.Msg)
	}
	return &cor.Data, nil
}

type SendMailRequest struct {
	From       string `json:"from" binding:"required"`
	To         string `json:"to" binding:"required"`
	HtmlString string `json:"html_string" binding:"required"`
}

type SendMailResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (c *UsdtPaymentClient) SendMail(ctx context.Context, req *SendMailRequest) error {
	r := resty.New().R().SetContext(ctx).SetBody(req).SetHeader("Authorization", c.token)
	rsp, err := r.Post(fmt.Sprintf("%s/api/v1/order/mail", c.host))
	if err != nil {
		return NewCodeError(500, err.Error())
	}
	var cor SendMailResponse
	err = json.Unmarshal(rsp.Body(), &cor)
	if err != nil {
		return err
	}
	if cor.Code != 0 {
		return NewCodeError(cor.Code, cor.Msg)
	}
	return nil
}
