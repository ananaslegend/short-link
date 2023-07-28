package main

import (
	"context"
	"github.com/ananaslegend/short-link/api/handlers/save"
	"github.com/ananaslegend/short-link/config"
	"github.com/ananaslegend/short-link/logs"
	"github.com/ananaslegend/short-link/storage"
	"golang.org/x/exp/slog"
	"net/http"

	"os"
)

func main() {
	confPath := os.Getenv("APP_CONFIG")
	cfg := config.MustLoadYaml(confPath)

	log := logs.SetUpLogger(cfg)
	log.Info("short-link app started", slog.String("env", string(cfg.Env)))

	db, err := storage.NewSqliteStorage(cfg.DbConn)
	if err != nil {
		log.Error("cant connect to database", logs.Err(err))
		os.Exit(1)
	}
	if err = db.PrepareStorage(); err != nil {
		log.Error("cant prepare database", logs.Err(err))
		os.Exit(1)
	}

	// setup router
	m := http.NewServeMux()

	m.HandleFunc("/link", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			save.Handle(context.TODO(), w, r, log, db)
		}
	})

	// start server
	s := http.Server{
		Addr:    cfg.HttpServer.Port,
		Handler: recoverHandler(log, m),
	}

	if err = s.ListenAndServe(); err != nil {
		log.Error("HTTP server", logs.Err(err))
		os.Exit(1)
	}
}

func recoverHandler(log *slog.Logger, m *http.ServeMux) http.Handler {
	defer func() {
		if r := recover(); r != nil {
			log.Error("app in panic", r)
		}
	}()

	return m
}
