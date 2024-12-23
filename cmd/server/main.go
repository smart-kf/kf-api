package main

import (
	"flag"
	xlogger "github.com/clearcodecn/log"
	"golang.org/x/sync/errgroup"
	"log"
	"std-api/config"
	"std-api/pkg/caches"
	"std-api/pkg/db"
	"std-api/pkg/server"
	"time"
)

var configName string

func init() {
	flag.StringVar(&configName, "c", "config.yaml", "配置文件")
}

func main() {
	flag.Parse()
	conf := config.Load(configName)
	initLogger(conf)
	db.Load()
	db.InitRedis()

	caches.InitCacheInstances()
	var (
		eg       errgroup.Group
		stopChan = make(chan struct{})
	)
	eg.Go(func() error {
		task := db.InitBillLogBackgroundTask(1*time.Minute, 100) // 1分钟清空buffer
		task.Start(stopChan)
		return nil
	})
	eg.Go(func() error {
		task := db.InitKFLogBackgroundTask(1*time.Minute, 10000)
		task.Start(stopChan)
		return nil
	})

	if err := server.Run(); err != nil {
		close(stopChan)
		return
	}

	log.Fatal(eg.Wait())
}

func initLogger(conf *config.Config) {
	//xlogger.AddHook(func(ctx context.Context) xlogger.Field {
	//	reqid, ok := ctx.Value("reqid").(string)
	//	if !ok {
	//		return xlogger.Field{}
	//	}
	//	return xlogger.Any("reqid", reqid)
	//})
	logger, err := xlogger.NewLog(xlogger.Config{
		Level:  conf.Log.Level,
		Format: conf.Log.Format,
		File:   conf.Log.File,
	})

	if err != nil {
		panic(err)
	}

	xlogger.SetGlobal(logger)
}
