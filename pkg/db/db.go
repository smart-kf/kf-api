package db

import (
	"context"
	xlogger "github.com/clearcodecn/log"
	"github.com/clearcodecn/sqlite"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"std-api/config"
	"sync"
)

var (
	o   sync.Once
	_db *gorm.DB
)

func Load() *gorm.DB {
	o.Do(func() {
		conf := config.GetConfig()
		var (
			db  *gorm.DB
			err error
		)
		switch conf.DB.Driver {
		case "sqlite":
			dir := filepath.Dir(conf.DB.Dsn)
			os.MkdirAll(dir, 0755)
			db, err = gorm.Open(sqlite.Open(conf.DB.Dsn), &gorm.Config{})
		case "mysql":
			db, err = gorm.Open(mysql.Open(conf.DB.Dsn), &gorm.Config{})
		default:
			log.Fatal("database driver not set")
		}

		if err != nil {
			log.Fatal("初始化db失败")
			return
		}
		_db = db

		db.Use(xlogger.NewLoggerPlugin())

		syncTable(_db)

		go initBillAccounts()
	})

	return _db
}

func syncTable(db *gorm.DB) {
	if err := db.AutoMigrate(
		&BillAccount{},
		&KFCard{},
		&BillLog{},
	); err != nil {
		log.Fatal(err)
	}
}

func DB() *gorm.DB {
	return _db
}

var transactionKey = struct{}{}

// Begin 开启事务
func Begin(ctx context.Context) (*gorm.DB, context.Context) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ctx = ginCtx.Request.Context()
	}
	db, ok := ctx.Value(transactionKey).(*gorm.DB)
	if !ok {
		db = DB()
		db = db.Begin()
		ctx = context.WithValue(ctx, transactionKey, db)
	}
	return db, ctx
}

func GetDBFromContext(ctx context.Context) *gorm.DB {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ctx = ginCtx.Request.Context()
	}
	db, ok := ctx.Value(transactionKey).(*gorm.DB)
	if !ok {
		db = DB()
		db = db.WithContext(ctx)
	}
	return db
}

func initBillAccounts() {
	conf := config.GetConfig()
	bac := conf.BillConfig.Accounts
	if len(bac) == 0 {
		return
	}

	for _, ac := range bac {
		acc := BillAccount{
			Username: ac.Username,
		}

		pass, _ := bcrypt.GenerateFromPassword([]byte(ac.Password), bcrypt.DefaultCost)
		acc.Password = string(pass)
		db := DB()
		var dbAcc BillAccount
		db.Where("username = ?", ac.Username).First(&dbAcc)

		if dbAcc.ID > 0 {
			if err := bcrypt.CompareHashAndPassword([]byte(dbAcc.Password), []byte(ac.Password)); err != nil {
				db.Model(BillAccount{}).Where("id = ?", dbAcc.ID).Update("password", acc.Password)
			}
		} else {
			db.Create(&acc)
		}
	}
}
