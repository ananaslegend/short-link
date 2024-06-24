package app

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ananaslegend/go-logs/v2"
	"github.com/ananaslegend/short-link/internal/config"
	"github.com/ananaslegend/short-link/internal/metrics"
	"github.com/ananaslegend/short-link/internal/middleware"
	redirectHandler "github.com/ananaslegend/short-link/internal/redirect/handler"
	redirectRepo "github.com/ananaslegend/short-link/internal/redirect/repository"
	redirectService "github.com/ananaslegend/short-link/internal/redirect/service"
	saveHandler "github.com/ananaslegend/short-link/internal/save/handler"
	saveRepo "github.com/ananaslegend/short-link/internal/save/repository"
	saveService "github.com/ananaslegend/short-link/internal/save/service"
	"github.com/ananaslegend/short-link/internal/statistic"
	"github.com/ananaslegend/short-link/internal/storage"
	"github.com/ananaslegend/short-link/pkg/clog"
	"github.com/go-pkgz/routegroup"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type StatManager interface {
	AppendRow(row *statistic.Row)
	Close(ctx context.Context) error
	Run()
}

type App struct {
	logger       *slog.Logger
	config       config.AppConfig
	db           *sql.DB
	redisClient  *redis.Client
	httpServer   *http.Server
	metricServer *http.Server
	statManager  StatManager
}

func New(ctx context.Context, confPath string) App {
	a := App{}

	a.config = config.MustLoadYaml(confPath)

	ctx = a.setUpLogger(ctx)

	a.metricServer = metrics.SetUpServer(a.config.Metrics.Addr)

	a.mustSetUpDB(ctx)

	a.setUpRedis()

	a.setUpStatisticManager()

	a.setUpHTTPServer()

	return a
}

func (a *App) Run() error {
	eg := errgroup.Group{} // todo WithContext

	eg.Go(func() error {
		a.logger.Info("start metric server", slog.String("addr", a.config.Metrics.Addr))
		if err := a.metricServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("metrics server", logs.ErrorMsg(err))
			return err
		}

		return nil
	})

	eg.Go(func() error {
		a.logger.Info("start HTTP server", slog.String("port", a.config.HttpServer.Port))
		if err := a.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("HTTP server", logs.ErrorMsg(err))
			return err
		}

		return nil
	})

	eg.Go(func() error {
		a.logger.Info("start statistic manager")
		a.statManager.Run()

		return nil
	})

	if err := eg.Wait(); err != nil {
		a.logger.Error("error during app.Run()", logs.ErrorMsg(err))

		return err
	}

	return nil
}

func (a *App) Close() {
	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		a.logger.Error("error during shutdown http server", logs.ErrorMsg(err))
	}

	if err := a.metricServer.Shutdown(context.Background()); err != nil {
		a.logger.Error("error during shutdown metric server", logs.ErrorMsg(err))
	}

	if err := a.statManager.Close(context.Background()); err != nil {
		a.logger.Error("error during close statistic manager", logs.ErrorMsg(err))
	}
}

func (a *App) setUpHTTPServer() {
	router := routegroup.New(http.NewServeMux())

	router.Use(clog.WithCtxLogger(a.logger))
	router.Use(middleware.WithRequestID)
	router.Use(middleware.WithRecover)
	router.Use(middleware.WithLoggingRequest)

	router.HandleFunc("GET /{alias}", a.redirectHandler())
	router.HandleFunc("POST /link", a.saveLinkHandler())

	a.httpServer = &http.Server{
		Addr:    a.config.HttpServer.Port,
		Handler: router,
	}
}

func (a *App) saveLinkHandler() http.HandlerFunc {
	repositorySave := saveRepo.NewPostgres(a.db)
	serviceSave := saveService.New(repositorySave)

	return saveHandler.New(serviceSave, a.logger).ServeHTTP
}

func (a *App) redirectHandler() http.HandlerFunc {
	repositoryRedirect := redirectRepo.New(a.db)

	cachedRepositoryRedirect := redirectRepo.NewRedisCache(repositoryRedirect, a.redisClient)

	serviceRedirect := redirectService.New(a.logger, cachedRepositoryRedirect, a.statManager)

	return redirectHandler.New(serviceRedirect).ServeHTTP
}

func (a *App) setUpLogger(ctx context.Context) context.Context {
	var handler slog.Handler

	switch a.config.Env {
	case config.Local:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case config.Dev:
		handler = slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	default:
		handler = slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true}))
	}

	logger := slog.New(handler)

	a.logger = logger

	return a.setAppLoggingMetrics(ctx)
}

func (a *App) setAppLoggingMetrics(ctx context.Context) context.Context {
	ctx = logs.WithMetric(ctx, "app", "short-link")
	ctx = logs.WithMetric(ctx, "env", a.config.Env)

	return ctx
}

func (a *App) mustSetUpDB(ctx context.Context) {
	db, err := storage.NewPostgres(a.config.DbConn)
	if err != nil {
		a.logger.ErrorContext(ctx, "cant connect to database", logs.ErrorMsg(err))
		os.Exit(1)
	}

	a.db = db

	a.logger.InfoContext(ctx, "database connected")
}

func (a *App) setUpRedis() {
	a.redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})
}

func (a *App) setUpStatisticManager() {
	repositoryStatistic := statistic.NewRepository(a.db)
	a.statManager = statistic.NewManager(1*time.Minute, 1000, repositoryStatistic, a.logger)
}
