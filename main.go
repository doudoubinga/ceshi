package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/doudoubinga/proxypool/api"
	"github.com/doudoubinga/proxypool/internal/app"
	"github.com/doudoubinga/proxypool/internal/cron"
	"github.com/doudoubinga/proxypool/internal/database"
	"github.com/doudoubinga/proxypool/pkg/proxy"
)

var configFilePath = ""

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	flag.StringVar(&configFilePath, "c", "", "path to config file: config.yaml")
	flag.Parse()

	if configFilePath == "" {
		configFilePath = os.Getenv("CONFIG_FILE")
	}
	if configFilePath == "" {
		configFilePath = "config.yaml"
	}
	err := app.InitConfigAndGetters(configFilePath)
	if err != nil {
		panic(err)
	}

	database.InitTables()
	proxy.InitGeoIpDB()
	fmt.Println("Do the first crawl...")
	go app.CrawlGo()
	go cron.Cron()
	api.Run()
}
