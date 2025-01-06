package main

import (
	"flag"
	"log"
	"time"

	xlogger "github.com/clearcodecn/log"
	"golang.org/x/sync/errgroup"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/endpoints/cron/billlog"
	"github.com/smart-fm/kf-api/endpoints/cron/kflog"
	"github.com/smart-fm/kf-api/endpoints/http"
	"github.com/smart-fm/kf-api/endpoints/nsq/producer"
	"github.com/smart-fm/kf-api/infrastructure/caches"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/nsq"
	"github.com/smart-fm/kf-api/infrastructure/redis"
)

var configName string

func init() {
	flag.StringVar(&configName, "c", "config.yaml", "配置文件")
}

func main() {
	flag.Parse()
	conf := config.Load(configName)
	initLogger(conf)
	mysql.Load()
	redis.InitRedis()
	nsq.InitNSQ()
	producer.InitProducer()

	caches.InitCacheInstances()
	var (
		eg       errgroup.Group
		stopChan = make(chan struct{})
	)
	eg.Go(
		func() error {
			task := billlog.InitBillLogBackgroundTask(1*time.Minute, 100) // 1分钟清空buffer
			task.Start(stopChan)
			return nil
		},
	)
	eg.Go(
		func() error {
			task := kflog.InitKFLogBackgroundTask(1*time.Minute, 10000)
			task.Start(stopChan)
			return nil
		},
	)
	eg.Go(
		func() error {
			nsq.StartConsume(stopChan)
			return nil
		},
	)

	if err := http.Run(); err != nil {
		close(stopChan)
		return
	}

	log.Fatal(eg.Wait())
}

func initLogger(conf *config.Config) {
	// xlogger.AddHook(func(ctx context.Context) xlogger.Field {
	//	reqid, ok := ctx.Value("reqid").(string)
	//	if !ok {
	//		return xlogger.Field{}
	//	}
	//	return xlogger.Any("reqid", reqid)
	// })
	logger, err := xlogger.NewLog(
		xlogger.Config{
			Level:  conf.Log.Level,
			Format: conf.Log.Format,
			File:   conf.Log.File,
		},
	)

	if err != nil {
		panic(err)
	}

	xlogger.SetGlobal(logger)
}
