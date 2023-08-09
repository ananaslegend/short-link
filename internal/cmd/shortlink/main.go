package main

import (
	handlers2 "github.com/ananaslegend/short-link/internal/api/handlers"
	"github.com/ananaslegend/short-link/internal/api/mw"
	"github.com/ananaslegend/short-link/internal/config"
	"github.com/ananaslegend/short-link/internal/services/link"
	"github.com/ananaslegend/short-link/internal/storage"
	"github.com/ananaslegend/short-link/pkg/logs"
	"log/slog"
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
	defer db.Close()

	linkService := link.New(db)

	m := http.NewServeMux()

	m.HandleFunc("/", mw.WithRequestId(
		func(w http.ResponseWriter, r *http.Request) {
			handlers2.Redirect(w, r, log, linkService)
		}))

	m.HandleFunc("/link", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers2.SaveLink(w, r, log, linkService)
		}
	})

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
