package main

import (
	"context"
	"errors"
	"flag"
	"github.com/ananaslegend/go-logs/v2"
	"github.com/ananaslegend/short-link/internal/config"
	"github.com/ananaslegend/short-link/internal/metrics"
	"github.com/ananaslegend/short-link/internal/middleware"
	redirectHandler "github.com/ananaslegend/short-link/internal/redirect/handler"
	redirectSqlite "github.com/ananaslegend/short-link/internal/redirect/repository"
	redirectService "github.com/ananaslegend/short-link/internal/redirect/service"
	saveHandler "github.com/ananaslegend/short-link/internal/save/handler"
	saveSqlite "github.com/ananaslegend/short-link/internal/save/repository"
	saveService "github.com/ananaslegend/short-link/internal/save/service"
	"github.com/ananaslegend/short-link/internal/statistic"
	"github.com/ananaslegend/short-link/internal/storage/cache"
	sql2 "github.com/ananaslegend/short-link/internal/storage/sqlutil"
	"github.com/ananaslegend/short-link/pkg/closer"
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

	logger := SetUpLogger(cfg)
	logger.Info("short-link app started", slog.String("env", string(cfg.Env)))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	gracefulCloser := closer.New()

	db, err := sql2.NewSQLiteStorage(cfg.DbConn)
	if err != nil {
		logger.Error("cant connect to database", logs.ErrorMsg(err))
		os.Exit(1)
	}
	logger.Debug("database connected")
	defer sql2.Close(db, logger)

	err = sql2.Prepare(db)
	if err != nil {
		logger.Error("cant prepare database", logs.ErrorMsg(err))
		os.Exit(1)
	}
	logger.Debug("database prepared")

	go func() {
		if err := metrics.Listen(cfg.Metrics.Addr); err != nil {
			logger.Error("cant listen metrics", logs.ErrorMsg(err))
		}
	}()
	logger.Info("metrics server started")

	linkCache, err := cache.NewCache(cfg.LinkCache)
	if err != nil {
		logger.Error("cant create link_cache", logs.ErrorMsg(err))
		os.Exit(1)
	}

	var (
		repositoryStatistic = statistic.NewRepository(db)
		statisticManager    = statistic.NewManager(1*time.Minute, 1000, repositoryStatistic, logger)
	)
	go statisticManager.Run()
	gracefulCloser.Add(statisticManager.Close)

	m := http.NewServeMux()

	var (
		repositoryRedirect       = redirectSqlite.New(db)
		cachedRepositoryRedirect = redirectSqlite.NewCached(repositoryRedirect, linkCache)
		serviceRedirect          = redirectService.New(logger, cachedRepositoryRedirect, statisticManager)
		handlerRedirect          = redirectHandler.New(serviceRedirect, logger)
	)
	m.HandleFunc("/", middleware.WithRequestID(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlerRedirect.ServeHTTP(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}),
	)

	var (
		repositorySave = saveSqlite.New(db)
		serviceSave    = saveService.New(logger, repositorySave)
		handlerSave    = saveHandler.New(serviceSave, logger)
	)
	m.HandleFunc("/link", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlerSave.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	s := http.Server{
		Addr:    cfg.HttpServer.Port,
		Handler: middleware.WithRecover(logger, m),
	}

	gracefulCloser.Add(s.Shutdown)

	go func() {
		if err = s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server", logs.ErrorMsg(err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	logger.Info("shutting down server gracefully")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutDownTimeout)
	defer cancel()

	if err = gracefulCloser.Close(shutdownCtx); err != nil {
		logger.Error("graceful shutdown with errors", logs.ErrorMsg(err))
	}
	logger.Info("graceful shutdown finished")
}

func SetUpLogger(cfg config.AppConfig) *slog.Logger { // todo move
	var logger *slog.Logger

	switch cfg.Env {
	case config.Local:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return logger
}
