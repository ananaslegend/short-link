package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Env string

const (
	Dev   Env = "dev"
	Prod  Env = "prod"
	Local Env = "local"
)

type AppConfig struct {
	Env        Env        `yaml:"env"`
	DbConn     string     `yaml:"db_conn" env-required:"true"`
	HttpServer HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Port string `yaml:"port"`
	// TODO timeouts
}

func MustLoadYaml(confPath string) AppConfig {
	if confPath == "" {
		log.Fatalf("config path is empty")
	}

	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", confPath)
	}

	cfg := &AppConfig{
		HttpServer: HttpServer{},
	}
	if data, err := os.ReadFile(confPath); err != nil {
		log.Fatalf("failed to read config file: %s", err)
	} else {
		if err = yaml.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("failed to unmarshal config file: %s", err)
		}
	}

	return *cfg
}
