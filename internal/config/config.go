package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type Env string

const (
	Dev   Env = "dev"
	Prod  Env = "prod"
	Local Env = "local"
)

type CacheType string

const (
	BigCache = "bigcache"
	Redis    = "redis"
)

type Cache struct {
	TTL       time.Duration `yaml:"ttl"`
	CacheType CacheType     `yaml:"type"`
}

type AppConfig struct {
	Env             Env           `yaml:"env"`
	DbConn          string        `yaml:"db_conn" env-required:"true"`
	HttpServer      HttpServer    `yaml:"http_server"`
	LinkCache       Cache         `yaml:"link_cache"`
	ShutDownTimeout time.Duration `yaml:"shut_down_timeout"`
	Metrics         Metrics       `yaml:"metrics"`
}

type Metrics struct {
	Addr string `json:"addr,omitempty"`
}

type HttpServer struct {
	Port string `yaml:"port"`
	// TODO timeouts
}

// MustLoadYaml loads config from yaml file and panic if error occurred.
// Config path have format: ../path/to/config.yaml
func MustLoadYaml(confPath string) AppConfig {
	if confPath == "" {
		log.Fatalf("config path is empty")
	}

	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", confPath)
	}

	cfg := AppConfig{
		HttpServer: HttpServer{},
	}
	dir, _ := os.Getwd()
	log.Printf("dir: %s", dir)
	if data, err := os.ReadFile(confPath); err != nil {
		log.Fatalf("failed to read config file: %s", err)
	} else {
		if err = yaml.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("failed to unmarshal config file: %s", err)
		}
	}

	return cfg
}
