package kffrontend

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kffrontend"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type QRCodeController struct {
	BaseController
}

func (c *QRCodeController) Scan(ctx *gin.Context) {
	var req kffrontend.QRCodeScanRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	return
}

func (c *BaseController) Upload(ctx *gin.Context) {
	var req kfbackend.UploadRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()

	kfToken := common.GetKFToken(reqCtx)
	cardId, err := caches.KfAuthCacheInstance.GetFrontToken(ctx.Request.Context(), kfToken)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	fi, err := req.File.Open()
	if err != nil {
		c.Error(ctx, err)
		return
	}
	defer fi.Close()

	head := make([]byte, 261)
	n, err := fi.Read(head)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	fi.Seek(0, io.SeekStart)

	if req.FileType == "image" && !filetype.IsImage(head[:n]) {
		c.Error(ctx, xerrors.NewCustomError("仅允许上传图片"))
		return
	}

	if req.FileType == "video" && !filetype.IsVideo(head[:n]) {
		c.Error(ctx, xerrors.NewCustomError("仅允许上传视频"))
		return
	}

	typ, err := filetype.Match(head[:n])
	if err != nil {
		c.Error(ctx, err)
		return
	}

	id := uuid.New()
	date := time.Now().Format(`200601`)
	// 1. 创建本地文件夹.
	fp := filepath.Join(config.GetConfig().Web.UploadDir, date, fmt.Sprintf("%s.%s", id, typ.Extension))
	if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		c.Error(ctx, xerrors.NewCustomError("创建bucket失败"))
		return
	}

	dst, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		c.Error(ctx, xerrors.NewCustomError("上传失败"))
		return
	}
	defer dst.Close()

	buffer := make([]byte, 1024*1024) // 1mb chunk
	// 创建 MD5 哈希对象
	hash := md5.New()
	for {
		n, err := fi.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if n == 0 {
			break
		}
		// 写入到磁盘
		if _, err := dst.Write(buffer[:n]); err != nil {
			c.Error(ctx, xerrors.NewCustomError("上传失败"))
			return
		}
		// 更新 MD5 哈希
		hash.Write(buffer[:n])
	}

	// 计算 MD5 值
	md5Sum := hash.Sum(nil)
	md5Hex := hex.EncodeToString(md5Sum)

	var file = dao.KfFile{
		Model:      gorm.Model{},
		Filename:   fmt.Sprintf("%s.%s", id, typ.Extension),
		Ext:        typ.Extension,
		FileType:   typ.MIME.Value,
		Md5:        md5Hex,
		CardId:     cardId,
		PublicPath: "/" + filepath.Join(date, fmt.Sprintf("%s.%s", id, typ.Extension)),
		UUID:       kfToken,
	}

	tx := mysql.DB()
	if err := tx.Create(&file).Error; err != nil {
		c.Error(ctx, xerrors.NewCustomError("上传失败"))
		return
	}

	c.Success(
		ctx, kfbackend.UploadResponse{
			Path:    file.PublicPath,
			CdnHost: config.GetConfig().Web.CdnHost,
		},
	)
}
