package repository

import (
	"context"

	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type KfUserRepository struct{}

func (KfUserRepository) CreateOne(ctx context.Context, user *dao.KfUsers) error {
	db := mysql.GetDBFromContext(ctx)
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
