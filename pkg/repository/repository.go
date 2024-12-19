package repository

import (
	"gorm.io/gorm"
	"std-api/pkg/db"
)

type BaseRepository struct{}

func NewBaseRepository() *BaseRepository {
	return &BaseRepository{}
}

func (b *BaseRepository) getDB() *gorm.DB {
	db := db.DB()
	return db
}
