package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string `validate:"required"`
	Environment Env    `validate:"required,oneof=prod dev local"`
	DbConn      string
	HttpServer  HttpServer
	Metrics     Metrics
	ClickHouse  ClickHouse
	Swagger     Swagger
	Redis       Redis
	Otel        Otel

	FlushStatisticDuration time.Duration
}

type Env string

const (
	Production = "prod"
	Dev        = "dev"
	Local      = "local"
)

type HttpServer struct {
	Port string `validate:"required"`
}

type Redis struct {
	Addr     string `validate:"required"`
	Password string `validate:"required"`
}

type Swagger struct {
	Port int `validate:"required"`
}

type ClickHouse struct {
	Host string `validate:"required"`
	Db   string `validate:"required"`
	User string `validate:"required"`
	Pass string `validate:"required"`
}

type Metrics struct {
	Addr string
}

type Otel struct {
	TraceGRCPAddr       string
	TraceFlushInterval  time.Duration
	MetricFlushInterval time.Duration
}

func MustLoadConfig() Config {
	_ = godotenv.Load(".env")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config

	cfg.Environment = Env(viper.GetString("ENV"))
	cfg.ServiceName = viper.GetString("SERVICE_NAME")

	cfg.HttpServer.Port = viper.GetString("HTTP_SERVER_PORT")

	cfg.DbConn = viper.GetString("DB_CONN")

	cfg.Metrics.Addr = os.Getenv("METRICS_ADDR")

	cfg.ClickHouse.Host = viper.GetString("CLICKHOUSE_HOST")
	cfg.ClickHouse.Db = viper.GetString("CLICKHOUSE_DATABASE")
	cfg.ClickHouse.User = viper.GetString("CLICKHOUSE_USER")
	cfg.ClickHouse.Pass = viper.GetString("CLICKHOUSE_PASSWORD")

	swaggerPort, err := strconv.Atoi(viper.GetString("SWAGGER_PORT"))
	if err != nil {
		panic(fmt.Sprintf("invalid swagger documentation http server port: %v", err))
	}

	cfg.Swagger.Port = swaggerPort

	cfg.Redis.Addr = viper.GetString("REDIS_ADDR")
	cfg.Redis.Password = viper.GetString("REDIS_PASS")

	cfg.FlushStatisticDuration = viper.GetDuration("FLUSH_STATISTIC_DURATION")

	cfg.Otel.TraceGRCPAddr = viper.GetString("OTEL_TRACE_GRCP_ADDR")
	cfg.Otel.TraceFlushInterval = viper.GetDuration("OTEL_TRACE_FLUSH_INTERVAL")
	cfg.Otel.MetricFlushInterval = viper.GetDuration("OTEL_METRIC_FLUSH_INTERVAL")

	if err = validator.New().Struct(cfg); err != nil {
		panic(fmt.Sprintf("config validation failed: %v", err))
	}

	return cfg
}
