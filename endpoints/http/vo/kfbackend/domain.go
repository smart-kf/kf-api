package kfbackend

type DomainOrderRequest struct {
	Num int64 `json:"num" binding:"required"`
}
