package main

import (
	"github.com/ananaslegend/short-link/config"
	"github.com/ananaslegend/short-link/logs"
	"github.com/ananaslegend/short-link/storage"
	"golang.org/x/exp/slog"

	"os"
)

func main() {
	confPath := os.Getenv("APP_CONFIG")
	cfg := config.MustLoadYaml(confPath)

	log := logs.SetUpLogger(cfg)
	log.Info("link-shorter app started", slog.String("env", string(cfg.Env)))

	// setup db
	db, err := storage.NewSqliteStorage(cfg.DbConn)
	log.Error("cant connect to database", logs.Err(err))
	_ = db
	// setup router

	// start server
}
