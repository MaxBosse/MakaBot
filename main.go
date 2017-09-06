package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/HeroesAwaken/GoAwaken/core"
	"github.com/MaxBosse/MakaBot/bot"
	"github.com/MaxBosse/MakaBot/log"
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

	metricConnection := new(core.InfluxDB)
	err = metricConnection.New(MyConfig.InfluxDBHost, MyConfig.InfluxDBDatabase, MyConfig.InfluxDBUser, MyConfig.InfluxDBPassword, AppName, Version)
	if err != nil {
		log.Fatalln("Error connecting to MetricsDB:", err)
	}

	bot.NewMakaBot(metricConnection, MyConfig.Servers, mem, MyConfig.DiscordToken)

	log.Noteln("MakaBot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
