//go:build sqlite
// +build sqlite

package db

import (
	"github.com/clearcodecn/sqlite"
	"github.com/smart-fm/kf-api/config"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

func init() {
	buildSqliteFunc = func() (*gorm.DB, error) {
		conf := config.GetConfig()
		dir := filepath.Dir(conf.DB.Dsn)
		os.MkdirAll(dir, 0755)
		return gorm.Open(sqlite.Open(conf.DB.Dsn), &gorm.Config{})
	}
}
