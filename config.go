package main

import (
	"io/ioutil"
	"log"

	"github.com/MaxBosse/MakaBot/bot/structs"
	"gopkg.in/yaml.v2"
)

type Config struct {
	InfluxDBHost     string
	InfluxDBDatabase string
	InfluxDBUser     string
	InfluxDBPassword string
	DiscordToken     string
	Servers          []*structs.DiscordServer
}

func (config *Config) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	return nil
}

func (config *Config) Load(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}
}
