package common

import (
	"errors"
	"github.com/modern-go/reflect2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	defaultScrollSize = 20
)

type ScrollRequest struct {
	Key      string      `json:"-"` // 基于哪个键滚页 必填
	ScrollID interface{} `json:"scrollID,omitempty" doc:"滚页的id"`
	PageSize *uint       `json:"pageSize,omitempty" doc:"分页大小,默认20"`
}

func (r *ScrollRequest) GetPageSize() int64 {
	if r.PageSize == nil {
		return defaultScrollSize
	}
	return int64(*r.PageSize)
}

func Scroll[T schema.Tabler](db *gorm.DB, request *ScrollRequest) ([]T, int64, error) {
	if len(request.Key) == 0 {
		return nil, 0, errors.New("key is required")
	}

	var cnt int64
	if err := db.Model(new(T)).Count(&cnt).Error; err != nil {
		return nil, 0, err
	}
	var res []T

	if reflect2.IsNil(request.ScrollID) {
		result := db.
			Order(request.Key + " asc").
			Limit(int(request.GetPageSize())).
			Find(&res)
		if result.Error != nil {
			return nil, 0, result.Error
		}
		return res, cnt, nil
	}

	result := db.
		Where(request.Key+" > ?", request.ScrollID).
		Order(request.Key + " asc").
		Limit(int(request.GetPageSize())).
		Find(&res)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return res, cnt, nil
}
