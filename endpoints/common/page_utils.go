package common

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
)

type PageRequest struct {
	Page     *uint  `json:"page" doc:"分页,从1开始"`       // 不要直接使用.
	PageSize *uint  `json:"pageSize" doc:"分页大小,默认20"` // 不要直接使用.
	OrderBy  string `json:"order_by" doc:"排序,默认id"`   // 排序字段: 默认传 id
	Asc      bool   `json:"asc" doc:"是否升序，默认倒序"`
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
	if request.OrderBy == "" {
		db = db.Order("id desc")
	} else {
		if request.Asc {
			db = db.Order(fmt.Sprintf("%s", request.OrderBy))
		} else {
			db = db.Order(fmt.Sprintf("%s desc", request.OrderBy))
		}
	}
	var res []T
	result := db.Offset(int((request.GetPage() - 1) * request.GetPageSize())).Limit(int(request.GetPageSize())).Find(&res)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return res, cnt, nil
}
