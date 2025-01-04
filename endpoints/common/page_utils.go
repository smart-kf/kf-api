package common

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
)

type PageRequest struct {
	Page     *uint `json:"page" doc:"分页,从1开始"`       // 不要直接使用.
	PageSize *uint `json:"pageSize" doc:"分页大小,默认20"` // 不要直接使用.
}

// GetPage 使用 GetPage 代替直接使用 .Page
func (r *PageRequest) GetPage() int64 {
	if r.Page == nil {
		return defaultPage
	}
	return int64(*r.Page)
}

func (r *PageRequest) GetPageSize() int64 {
	if r.PageSize == nil {
		return defaultPageSize
	}
	return int64(*r.PageSize)
}

func Paginate[T schema.Tabler](db *gorm.DB, request *PageRequest) ([]T, int64, error) {
	var cnt int64
	if err := db.Model(new(T)).Count(&cnt).Error; err != nil {
		return nil, 0, err
	}
	var res []T
	result := db.Offset(int((request.GetPage() - 1) * request.GetPageSize())).Limit(int(request.GetPageSize())).Find(&res)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return res, cnt, nil
}
