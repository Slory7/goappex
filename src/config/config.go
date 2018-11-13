package config

import (
	"github.com/nuveo/log"

	"github.com/crgimenes/goconfig"
)

type Config struct {
	AppIsDebug                 bool     `json:"appIsDebug" cfg:"appIsDebug"`
	Addr                       string   `json:"addr" cfg:"addr"`
	DBType                     string   `json:"dbType" cfg:"dbType"`
	DBConnectionString         string   `json:"dbConnectionString" cfg:"dbConnectionString"`
	DBReadOnlyConnectionString string   `json:"dbReadOnlyConnectionString" cfg:"dbReadOnlyConnectionString"`
	Redis                      RedisCfg `json:"redis" cfg:"redis"`
}

type RedisCfg struct {
	Hosts        string `json:"hosts" cfg:"hosts"`
	Password     string `json:"password" cfg:"password"`
	MasterName   string `json:"masterName" cfg:"masterName"`
	DBNumber     int    `json:"dbNumber" cfg:"dbNumber"`
	MaxRetries   int    `json:"maxRetries" cfg:"maxRetries" cfgDefault:"0"`
	PoolSize     int    `json:"poolSize" cfg:"poolSize" cfgDefault:"1000"`
	IdleTimeout  int    `json:"idleTimeout" cfg:"idleTimeout" cfgDefault:"300"`
	MinIdleConns int    `json:"minIdleConns" cfg:"minIdleConns" cfgDefault:"0"`
}

func GetConfig(environment string) *Config {
	config := Config{}
	if environment != "" {
		goconfig.File = "config." + environment + ".json"
	} else {
		goconfig.File = "config.json"
	}
	err := goconfig.Parse(&config)
	if err != nil {
		log.Fatal("config: ", err)
	}
	return &config
}
