package main

import (
	"github.com/ananaslegend/short-link/internal/middleware"
	"github.com/ananaslegend/short-link/internal/redirect"
	"github.com/ananaslegend/short-link/internal/save"
	"github.com/ananaslegend/short-link/pkg/config"
	"github.com/ananaslegend/short-link/pkg/logs"
	"github.com/ananaslegend/short-link/pkg/storage/cache"
	"github.com/ananaslegend/short-link/pkg/storage/sql"
	"log/slog"
	"net/http"

	"os"
)

func main() {
	confPath := os.Getenv("APP_CONFIG")
	cfg := config.MustLoadYaml(confPath)

	log := logs.SetUpLogger(cfg)
	log.Info("short-link app started", slog.String("env", string(cfg.Env)))

	db, err := sql.NewSqliteStorage(cfg.DbConn)
	if err != nil {
		log.Error("cant connect to database", logs.Err(err))
		os.Exit(1)
	}
	log.Debug("database connected")
	defer sql.Close(db, log)

	err = sql.Prepare(db)
	if err != nil {
		log.Error("cant prepare database", logs.Err(err))
		os.Exit(1)
	}
	log.Debug("database prepared")

	linkCache, err := cache.NewCache(cfg.LinkCache)
	if err != nil {
		log.Error("cant create link_cache", logs.Err(err))
		os.Exit(1)
	}

	m := http.NewServeMux()

	redirectRepo := redirect.NewRepository(db)
	cachedRedirectRepo := redirect.NewCachedRepository(redirectRepo, linkCache)
	redirectUseCase := redirect.NewUseCase(cachedRedirectRepo)
	redirectHandler := redirect.NewHandler(redirectUseCase, log)

	m.HandleFunc("/", middleware.WithRequestId(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				redirectHandler.HandleHTTP(w, r)
			}
		}),
	)

	saveRepo := save.NewRepository(db)
	saveUseCase := save.NewUseCase(saveRepo)
	saveHandler := save.NewHandler(saveUseCase, log)

	m.HandleFunc("/link", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			saveHandler.HandleHTTP(w, r)
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
