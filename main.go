package main

import (
	"flag"
	"github.com/guobang-yoo/proxypool/pkg/geoIp"
	_ "net/http/pprof"
	"os"

	"github.com/guobang-yoo/proxypool/api"
	"github.com/guobang-yoo/proxypool/internal/app"
	"github.com/guobang-yoo/proxypool/internal/cron"
	"github.com/guobang-yoo/proxypool/internal/database"
	"github.com/guobang-yoo/proxypool/log"
)

var configFilePath = ""
var debugMode = false

func main() {
	//go func() {
	//	http.ListenAndServe("0.0.0.0:6060", nil)
	//}()

	flag.StringVar(&configFilePath, "c", "", "path to config file: config.yaml")
	flag.BoolVar(&debugMode, "d", false, "debug output")
	flag.Parse()

	log.SetLevel(log.INFO)
	if debugMode {
		log.SetLevel(log.DEBUG)
	}
	if configFilePath == "" {
		configFilePath = os.Getenv("CONFIG_FILE")
	}
	if configFilePath == "" {
		configFilePath = "config.yaml"
	}
	err := app.InitConfigAndGetters(configFilePath)
	if err != nil {
		log.Errorln("Configuration init error: %s", err.Error())
		panic(err)
	}

	database.InitTables()
	// init GeoIp db reader and map between emoji's and countries
	// return: struct geoIp (dbreader, emojimap)
	err = geoIp.InitGeoIpDB()
	if err != nil {
		os.Exit(1)
	}
	log.Infoln("Do the first crawl...")
	go app.CrawlGo() // 抓取主程序
	go cron.Cron()   // 定时运行
	api.Run()        // Web Serve
}
