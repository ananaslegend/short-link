package main

import (
	"context"
	"errors"
	"flag"
	"github.com/ananaslegend/short-link/internal/metrics"
	"github.com/ananaslegend/short-link/internal/middleware"
	redirectHandler "github.com/ananaslegend/short-link/internal/redirect/handler"
	redirectSqlite "github.com/ananaslegend/short-link/internal/redirect/repository"
	redirectService "github.com/ananaslegend/short-link/internal/redirect/service"
	saveHandler "github.com/ananaslegend/short-link/internal/save/handler"
	saveSqlite "github.com/ananaslegend/short-link/internal/save/repository"
	saveService "github.com/ananaslegend/short-link/internal/save/service"
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

	go func() {
		if err := metrics.Listen(cfg.Metrics.Addr); err != nil {
			log.Error("cant listen metrics", logs.Err(err))
		}
	}()
	log.Info("metrics server started")

	linkCache, err := cache.NewCache(cfg.LinkCache)
	if err != nil {
		log.Error("cant create link_cache", logs.Err(err))
		os.Exit(1)
	}

	var (
		repositoryStatistic = statistic.NewRepository(db)
		statisticManager    = statistic.NewManager(1*time.Minute, 1000, repositoryStatistic, log)
	)
	go statisticManager.Run()
	gracefulCloser.Add(statisticManager.Close)

	m := http.NewServeMux()

	var (
		repositoryRedirect       = redirectSqlite.New(db)
		cachedRepositoryRedirect = redirectSqlite.NewCached(repositoryRedirect, linkCache)
		serviceRedirect          = redirectService.New(log, cachedRepositoryRedirect, statisticManager)
		handlerRedirect          = redirectHandler.New(serviceRedirect, log)
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
		serviceSave    = saveService.New(log, repositorySave)
		handlerSave    = saveHandler.New(serviceSave, log)
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
