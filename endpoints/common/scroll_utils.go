package common

import (
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
	Asc      bool        `json:"asc,omitempty" doc:"是否是升序"`
	ScrollID interface{} `json:"scrollID,omitempty" doc:"滚页的id"`
	PageSize *uint       `json:"pageSize,omitempty" doc:"分页大小,默认20"`
	Sorters  []Sorter    `json:"sorters" doc:"排序列表"`
}

type Sorter struct {
	Key string
	Asc bool
}

func (r *ScrollRequest) GetPageSize() int64 {
	if r.PageSize == nil {
		return defaultScrollSize
	}
	return int64(*r.PageSize)
}

func Scroll[T schema.Tabler](db *gorm.DB, request *ScrollRequest) ([]T, error) {
	var res []T

	var orderCols []clause.OrderByColumn

	if len(request.Sorters) == 0 {
		orderCols = append(
			orderCols, clause.OrderByColumn{
				Column: clause.Column{Name: "id"},
				Desc:   !request.Asc,
			},
		)
	} else {
		for _, sort := range request.Sorters {
			orderCols = append(
				orderCols,
				clause.OrderByColumn{Column: clause.Column{Name: sort.Key}, Desc: !sort.Asc},
			)
		}
	}
	order := clause.OrderBy{Columns: orderCols}
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
