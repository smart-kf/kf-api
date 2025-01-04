package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/infrastructure/mysql"
)

type BaseRepository struct{}

func NewBaseRepository() *BaseRepository {
	return &BaseRepository{}
}

func (b *BaseRepository) getDB(ctx context.Context) *gorm.DB {
	db := mysql.DB()
	return db.WithContext(ctx)
}
