package main

import (
	"context"
	"errors"
	"flag"
	"github.com/ananaslegend/short-link/internal/middleware"
	"github.com/ananaslegend/short-link/internal/redirect"
	"github.com/ananaslegend/short-link/internal/save"
	"github.com/ananaslegend/short-link/internal/statistic"
	"github.com/ananaslegend/short-link/pkg/closer"
	"github.com/ananaslegend/short-link/pkg/config"
	"github.com/ananaslegend/short-link/pkg/logs"
	"github.com/ananaslegend/short-link/pkg/storage/cache"
	"github.com/ananaslegend/short-link/pkg/storage/sql"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"os"
)

func main() {
	confPath := flag.String("config", "../../config/app-config.yml", "path to config file")
	flag.Parse()
	cfg := config.MustLoadYaml(*confPath)

	log := logs.SetUpLogger(cfg)
	log.Info("short-link app started", slog.String("env", string(cfg.Env)))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	gracefulCloser := closer.New()

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

	statRepo := statistic.NewRepository(db)
	statManager := statistic.NewManager(1*time.Minute, 1000, statRepo, log)
	go statManager.Run()
	gracefulCloser.Add(statManager.Close)

	m := http.NewServeMux()

	redirectRepo := redirect.NewRepository(db)
	cachedRedirectRepo := redirect.NewCachedRepository(redirectRepo, linkCache)
	redirectUseCase := redirect.NewUseCase(cachedRedirectRepo, statManager)
	redirectHandler := redirect.NewHandler(redirectUseCase, log)

	m.HandleFunc("/", middleware.WithRequestId(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				redirectHandler.ServeHTTP(w, r)
			}
		}),
	)

	saveRepo := save.NewRepository(db)
	saveUseCase := save.NewUseCase(saveRepo)
	saveHandler := save.NewHandler(saveUseCase, log)

	m.HandleFunc("/link", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			saveHandler.ServeHTTP(w, r)
		}
	})

	s := http.Server{
		Addr:    cfg.HttpServer.Port,
		Handler: middleware.WithRecover(log, m),
	}

	gracefulCloser.Add(s.Shutdown)

	go func() {
		if err = s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server", logs.Err(err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	log.Info("shutting down server gracefully")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutDownTimeout)
	defer cancel()

	if err = gracefulCloser.Close(shutdownCtx); err != nil {
		log.Error("graceful shutdown with errors", logs.Err(err))
	}
	log.Info("graceful shutdown finished")
}
