package repository

import (
	"context"
	"github.com/smart-fm/kf-api/pkg/db"
	"gorm.io/gorm"
)

type BaseRepository struct{}

func NewBaseRepository() *BaseRepository {
	return &BaseRepository{}
}

func (b *BaseRepository) getDB(ctx context.Context) *gorm.DB {
	db := db.DB()
	return db.WithContext(ctx)
}
