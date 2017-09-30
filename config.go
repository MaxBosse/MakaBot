package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Influx struct {
		Host     string
		Database string
		User     string
		Password string
	}
	MySQL struct {
		Host     string
		Database string
		User     string
		Password string
	}
	DiscordToken string
	WebBind      string
	CookieSecret string
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
