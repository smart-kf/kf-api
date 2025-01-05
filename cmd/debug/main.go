package main

import (
	"flag"

	xlogger "github.com/clearcodecn/log"

	"github.com/smart-fm/kf-api/config"
)

var configName string

func init() {
	flag.StringVar(&configName, "c", "config.yaml", "配置文件")
}

func main() {
	flag.Parse()
	conf := config.Load(configName)
	initLogger(conf)
	//
	// if err := server.Run(); err != nil {
	// 	// close(stopChan)
	// 	return
	// }

	// log.Fatal(eg.Wait())
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
