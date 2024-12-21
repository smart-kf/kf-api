//go:build sqlite
// +build sqlite

package db

import "gorm.io/gorm"

func init() {
	buildSqliteFunc = func() (*gorm.DB, error) {
		dir := filepath.Dir(conf.DB.Dsn)
		os.MkdirAll(dir, 0755)
		return gorm.Open(sqlite.Open(conf.DB.Dsn), &gorm.Config{})
	}
}
