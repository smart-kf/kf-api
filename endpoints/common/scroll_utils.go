package common

import (
	"errors"
	"github.com/modern-go/reflect2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

const (
	defaultScrollSize = 20
)

type ScrollRequest struct {
	Key      string      `json:"-"` // 基于哪个键滚页 必填
	Asc      bool        `json:"-"`
	ScrollID interface{} `json:"scrollID,omitempty" doc:"滚页的id"`
	PageSize *uint       `json:"pageSize,omitempty" doc:"分页大小,默认20"`
}

func (r *ScrollRequest) GetPageSize() int64 {
	if r.PageSize == nil {
		return defaultScrollSize
	}
	return int64(*r.PageSize)
}

func Scroll[T schema.Tabler](db *gorm.DB, request *ScrollRequest) ([]T, error) {
	if len(request.Key) == 0 {
		return nil, errors.New("key is required")
	}

	var res []T

	order := clause.OrderBy{Columns: []clause.OrderByColumn{
		{Column: clause.Column{Name: request.Key}, Desc: !request.Asc},
		{Column: clause.Column{Name: "id"}, Desc: false}, // 组合排序
	}}

	if reflect2.IsNil(request.ScrollID) {
		result := db.
			Order(order).
			Limit(int(request.GetPageSize())).
			Find(&res)
		if result.Error != nil {
			return nil, result.Error
		}
		return res, nil
	}

	result := db.
		Where(request.Key+" > ?", request.ScrollID).
		Order(order).
		Limit(int(request.GetPageSize())).
		Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return res, nil
}
