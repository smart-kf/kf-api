package kfbackend

import "mime/multipart"

type GetQRCodeIDRequest struct{}

type GetQRCodeIDResponse struct {
	Id string `json:"id"`
}

type UploadRequest struct {
	FileType string                `json:"fileType" form:"fileType" binding:"required" doc:"image || video 2选一"`
	File     *multipart.FileHeader `json:"file" form:"file" binding:"required" doc:"文件"`
}

type UploadResponse struct {
	Path    string `json:"path"`
	CdnHost string `json:"cdnHost"`
}
