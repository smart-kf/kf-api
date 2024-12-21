//go:build sqlite
// +build sqlite

package db

import (
	"github.com/clearcodecn/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"std-api/config"
)

func init() {
	buildSqliteFunc = func() (*gorm.DB, error) {
		conf := config.GetConfig()
		dir := filepath.Dir(conf.DB.Dsn)
		os.MkdirAll(dir, 0755)
		return gorm.Open(sqlite.Open(conf.DB.Dsn), &gorm.Config{})
	}
}
