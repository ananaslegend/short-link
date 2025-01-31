package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-pkgz/routegroup"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/errgroup"

	"github.com/ananaslegend/short-link/internal/config"
	"github.com/ananaslegend/short-link/internal/metrics"
	redirectHandler "github.com/ananaslegend/short-link/internal/redirect/handler"
	redirectRepo "github.com/ananaslegend/short-link/internal/redirect/repository"
	redirectService "github.com/ananaslegend/short-link/internal/redirect/service"
	saveHandler "github.com/ananaslegend/short-link/internal/save/handler"
	saveRepo "github.com/ananaslegend/short-link/internal/save/repository"
	saveService "github.com/ananaslegend/short-link/internal/save/service"
	"github.com/ananaslegend/short-link/internal/statistic"
	"github.com/ananaslegend/short-link/internal/storage"
	"github.com/ananaslegend/short-link/pkg/cslog"
)

const (
	defaultRowsCap = 1000
)

type StatManager interface {
	AppendRow(row *statistic.Row)
	Close(ctx context.Context) error
	Run()
	FlushTime() time.Duration
}

type App struct {
	logger        *slog.Logger
	config        config.AppConfig
	db            *sql.DB
	redisClient   *redis.Client
	httpServer    *http.Server
	metricServer  *http.Server
	swaggerServer *http.Server
	statManager   StatManager
}

func New(ctx context.Context) App {
	a := App{}

	a.config = config.MustLoadFromEnv()

	a.setupLogger()

	a.metricServer = metrics.SetupServer(a.config.Metrics.Addr)

	a.mustSetupPostgresDB(ctx)

	a.setupRedis()

	a.mustSetupStatisticManager()

	a.setupHTTPServer()

	a.setupSwaggerDocumentationServer()

	return a
}

func (a *App) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		a.logger.InfoContext(ctx, "start metric server", slog.String("addr", a.config.Metrics.Addr))
		if err := a.metricServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {

			return fmt.Errorf("failed to start metric server, error: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		a.logger.InfoContext(ctx, "start HTTP server", slog.String("port", a.config.HttpServer.Port))

		if err := a.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to start HTTP server, error: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		a.logger.InfoContext(ctx, "start statistic manager")
		a.statManager.Run()

		return nil
	})

	eg.Go(func() error {
		return a.runSwaggerDocumentationServer(ctx)
	})

	eg.Go(func() error {
		if res := a.redisClient.Ping(ctx); res != nil && res.Err() != nil {
			return fmt.Errorf("failed to ping redis, err: %w", res.Err())
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		a.logger.ErrorContext(ctx, "error during app.Run()", cslog.Error(err))

		return err
	}

	return nil
}

func (a *App) Close() {
	ctx, cancel := a.closeContext()
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.ErrorContext(ctx, "error during shutdown http server", cslog.Error(err))
	}

	if err := a.metricServer.Shutdown(ctx); err != nil {
		a.logger.ErrorContext(ctx, "error during shutdown metric server", cslog.Error(err))
	}

	if err := a.swaggerServer.Shutdown(ctx); err != nil {
		a.logger.ErrorContext(ctx, "error during close statistic manager", cslog.Error(err))
	}

	if res := a.redisClient.Shutdown(ctx); res != nil {
		a.logger.ErrorContext(ctx, "error during close redis client", cslog.Error(res.Err()))
	}

	if err := a.db.Close(); err != nil {
		a.logger.ErrorContext(ctx, "error during close db client", cslog.Error(err))
	}

	if err := a.statManager.Close(ctx); err != nil {
		a.logger.ErrorContext(ctx, "error during close statistic manager", cslog.Error(err))
	}
}

func (a *App) closeContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), a.config.ShutdownDuration)
}

func (a *App) setupHTTPServer() {
	router := routegroup.New(http.NewServeMux())

	a.setUpRouter(router)

	a.httpServer = &http.Server{
		Addr:    a.config.HttpServer.Port,
		Handler: router,

		ReadHeaderTimeout: config.DefaultReadHeaderRequestTimeout,
		ReadTimeout:       config.DefaultReadRequestTimeout,
		WriteTimeout:      config.DefaultWriteTimeout,
		IdleTimeout:       config.DefaultIdleTimeout,
	}
}

func (a *App) saveLinkHandler() http.HandlerFunc {
	repositorySave := saveRepo.NewPostgres(a.db)
	serviceSave := saveService.New(repositorySave)

	return saveHandler.New(serviceSave).ServeHTTP
}

func (a *App) redirectHandler() http.HandlerFunc {
	repositoryRedirect := redirectRepo.New(a.db)

	cachedRepositoryRedirect := redirectRepo.NewRedisCache(repositoryRedirect, a.redisClient)

	serviceRedirect := redirectService.New(cachedRepositoryRedirect, a.statManager)

	return redirectHandler.New(serviceRedirect).RedirectHandler
}

func (a *App) setupLogger() {
	var handler slog.Handler

	switch a.config.Env {
	case config.Local:
		handler = cslog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.Test:
		handler = cslog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	default:
		handler = cslog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true}))
	}

	logger := slog.New(handler)

	a.logger = logger.With(
		slog.String("service", "short-link"),
		slog.String("env", string(a.config.Env)),
	)
}

func (a *App) mustSetupPostgresDB(ctx context.Context) {
	db, err := storage.NewPostgres(a.config.DbConn)
	if err != nil {
		a.logger.ErrorContext(ctx, "cant connect to postgres database", cslog.Error(err))

		os.Exit(1)
	}

	a.db = db

	a.logger.InfoContext(ctx, "postgres database connected")
}

func (a *App) setupRedis() {
	a.redisClient = redis.NewClient(&redis.Options{
		Addr:     a.config.Redis.Addr,
		Password: a.config.Redis.Password,
	})
}

func (a *App) mustSetupStatisticManager() {
	ch, err := storage.NewClickHouse(a.config.ClickHouse.Host, a.config.ClickHouse.Port, a.config.ClickHouse.Db, a.config.ClickHouse.Pass, a.config.ClickHouse.User)
	if err != nil {
		a.logger.Error("cant connect to clickhouse", cslog.Error(err))
		os.Exit(1)
	}

	repositoryStatistic := statistic.NewNativeClickHouseRepository(ch)
	a.statManager = statistic.NewManager(1*time.Second, defaultRowsCap, repositoryStatistic, a.logger)
}
