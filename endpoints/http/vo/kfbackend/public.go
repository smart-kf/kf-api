package kfbackend

type GetQRCodeIDRequest struct{}

type GetQRCodeIDResponse struct {
	Id string `json:"id"`
}
