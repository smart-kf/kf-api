package kffrontend

type QRCodeScanRequest struct {
	Code string `json:"code" doc:"二维码的code" validate:"required" binding:"required"` // 二维码path
}
