package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/MaxBosse/MakaBot/bot"
	"github.com/MaxBosse/MakaBot/cache"
	"github.com/MaxBosse/MakaBot/log"
	"github.com/MaxBosse/MakaBot/utils"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	flag.StringVar(&configPath, "config", "config.yml", "Path to yml configuration file")
	flag.StringVar(&logLevel, "logLevel", "error", "LogLevel [error|warning|note|debug]")
	flag.Parse()

	log.SetLevel(logLevel)
	MyConfig.Load(configPath)
}

var (
	configPath string
	logLevel   string

	// MyConfig Default configuration
	MyConfig = Config{}

	Version = "0.0.1"
	mem     runtime.MemStats
	AppName = "MakaBot"
)

func main() {
	var err error

	if MyConfig.DiscordToken == "" {
		log.Errorln("No token provided. Please run: airhorn -t <bot token>")
		return
	}

	metricConnection := new(utils.InfluxDB)
	err = metricConnection.New(MyConfig.Influx.Host, MyConfig.Influx.Database, MyConfig.Influx.User, MyConfig.Influx.Password, AppName, Version)
	if err != nil {
		log.Fatalln("Error connecting to MetricsDB:", err)
	}

	// Collect memory statistics
	collectGlobalMetrics(metricConnection)
	globalTicker := time.NewTicker(time.Second * 10)
	go func() {
		for range globalTicker.C {
			collectGlobalMetrics(metricConnection)
		}
	}()

	dbConnection := new(utils.DB)
	dbSQL, err := dbConnection.New(MyConfig.MySQL.Host, MyConfig.MySQL.Database, MyConfig.MySQL.User, MyConfig.MySQL.Password)
	if err != nil {
		log.Fatalln("Error connecting to DB:", err)
	}

	cache := cache.NewCache(metricConnection, dbSQL)

	bot.NewMakaBot(MyConfig.DiscordToken, metricConnection, dbSQL, cache)

	log.Noteln("MakaBot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// CollectGlobalMetrics collects global metrics about the bot and environment
// And sends them to influxdb
func collectGlobalMetrics(iDB *utils.InfluxDB) {
	runtime.ReadMemStats(&mem)
	tags := map[string]string{"metric": "server_metrics", "server": "global", "serverID": "-1"}
	fields := map[string]interface{}{
		"memAlloc":      int(mem.Alloc),
		"memTotalAlloc": int(mem.TotalAlloc),
		"memHeapAlloc":  int(mem.HeapAlloc),
		"memHeapSys":    int(mem.HeapSys),
	}

	err := iDB.AddMetric("server_metrics", tags, fields)
	if err != nil {
		log.Errorln("Error adding Metric:", err)
	}
}
