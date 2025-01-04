//go:build sqlite
// +build sqlite

package dao

import (
	"os"
	"path/filepath"

	"github.com/clearcodecn/sqlite"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/infrastructure/mysql"

	"gorm.io/gorm"
)

func init() {
	mysql.buildSqliteFunc = func() (*gorm.DB, error) {
		conf := config.GetConfig()
		dir := filepath.Dir(conf.DB.Dsn)
		os.MkdirAll(dir, 0755)
		return gorm.Open(sqlite.Open(conf.DB.Dsn), &gorm.Config{})
	}
}
