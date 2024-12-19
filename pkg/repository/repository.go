package repository

import (
	"context"
	"gorm.io/gorm"
	"std-api/pkg/db"
)

type BaseRepository struct{}

func NewBaseRepository() *BaseRepository {
	return &BaseRepository{}
}

func (b *BaseRepository) getDB(ctx context.Context) *gorm.DB {
	db := db.DB()
	return db.WithContext(ctx)
}
