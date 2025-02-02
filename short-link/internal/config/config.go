package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	DefaultReadHeaderRequestTimeout = 5 * time.Second
	DefaultReadRequestTimeout       = 10 * time.Second
	DefaultWriteTimeout             = 10 * time.Second
	DefaultIdleTimeout              = 120 * time.Second
)

type Env string

const (
	Test  Env = "test"
	Prod  Env = "prod"
	Local Env = "local"
)

type CacheType string

type Cache struct {
	TTL       time.Duration `yaml:"ttl"`
	CacheType CacheType     `yaml:"type"`
}

type AppConfig struct {
	Env        Env        `env-required:"true"`
	DbConn     string     `env-required:"true"`
	HttpServer HttpServer `env-required:"true"`
	Metrics    Metrics    `env-required:"true"`
	ClickHouse ClickHouse `env-required:"true"`
	Swagger    Swagger    `env-required:"true"`
	Redis      Redis

	ShutdownDuration       time.Duration `env-required:"true"`
	FlushStatisticDuration time.Duration `env-required:"true"`
}

type Swagger struct {
	Port int
}

type ClickHouse struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Db   string `yaml:"db"`
	User string `yaml:"user"`
	Pass string `yaml:"password"`
}

type Metrics struct {
	Addr string `json:"addr,omitempty"`
}

type HttpServer struct {
	Port string `yaml:"port"`
	// TODO timeouts
}

type Redis struct {
	Addr     string
	Password string
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
	if data, err := os.ReadFile(confPath); err != nil {
		log.Fatalf("failed to read config file: %s", err)
	} else {
		if err = yaml.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("failed to unmarshal config file: %s", err)
		}
	}

	return cfg
}

func MustLoadFromEnv() AppConfig {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	if err := godotenv.Load(fmt.Sprintf(".env.%v", env)); err != nil {
		panic(fmt.Sprintf("failed to load env, error: %v", err))
	}

	cfg := AppConfig{}

	cfg.Env = Env(env)

	cfg.DbConn = os.Getenv("DB_CONN")

	cfg.HttpServer.Port = os.Getenv("HTTP_SERVER_PORT")

	cfg.Metrics.Addr = os.Getenv("METRICS_ADDR")

	cfg.ClickHouse.Host = os.Getenv("CLICKHOUSE_HOST")
	cfg.ClickHouse.Port = os.Getenv("CLICKHOUSE_PORT")
	cfg.ClickHouse.Db = os.Getenv("CLICKHOUSE_DATABASE")
	cfg.ClickHouse.User = os.Getenv("CLICKHOUSE_USER")
	cfg.ClickHouse.Pass = os.Getenv("CLICKHOUSE_PASSWORD")

	swaggerPort, err := strconv.Atoi(os.Getenv("SWAGGER_PORT"))
	if err != nil {
		panic(fmt.Sprintf("invalid swagger documentation http server port: %v", err))
	}

	cfg.Swagger.Port = swaggerPort

	shutdownDuration, err := time.ParseDuration(os.Getenv("SHUTDOWN_DURATION"))
	if err != nil {
		panic(fmt.Sprintf("invalid shutdown duration: %v", err))
	}

	cfg.ShutdownDuration = shutdownDuration

	cfg.Redis.Addr = os.Getenv("REDIS_ADDR")
	cfg.Redis.Password = os.Getenv("REDIS_PASS")

	statFlushDuration, err := time.ParseDuration(os.Getenv("FLUSH_STATISTIC_DURATION"))
	if err != nil || statFlushDuration == 0 {
		panic(fmt.Sprintf("invalid shutdown duration: %v", err))
	}

	cfg.FlushStatisticDuration = statFlushDuration

	return cfg
}
