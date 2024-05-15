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
	redirectSqlite "github.com/ananaslegend/short-link/internal/redirect/repository"
	redirectService "github.com/ananaslegend/short-link/internal/redirect/service"
	saveHandler "github.com/ananaslegend/short-link/internal/save/handler"
	saveSqlite "github.com/ananaslegend/short-link/internal/save/repository"
	saveService "github.com/ananaslegend/short-link/internal/save/service"
	"github.com/ananaslegend/short-link/internal/statistic"
	"github.com/ananaslegend/short-link/internal/storage/sqlutil"
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

	a.setUpStatisticManager()

	a.setUpHTTPServer()

	return a
}

func (a *App) Run() error {
	eg := errgroup.Group{} // todo WithContext

	eg.Go(func() error {
		if err := a.metricServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("metrics server", logs.ErrorMsg(err))
			return err
		}

		return nil
	})

	eg.Go(func() error {
		if err := a.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("HTTP server", logs.ErrorMsg(err))
			return err
		}

		return nil
	})

	eg.Go(func() error {
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
	router := http.NewServeMux()

	router.HandleFunc("GET /{alias}",
		middleware.WithRequestID(
			a.setUpRedirectHandler(),
		),
	)

	router.HandleFunc("POST /link", a.setUpSaveLinkHandler())

	a.httpServer = &http.Server{
		Addr:    a.config.HttpServer.Port,
		Handler: middleware.WithRecover(a.logger, router),
	}
}

func (a *App) setUpSaveLinkHandler() http.HandlerFunc {
	repositorySave := saveSqlite.New(a.db)
	serviceSave := saveService.New(a.logger, repositorySave)

	return saveHandler.New(serviceSave, a.logger).ServeHTTP
}

func (a *App) setUpRedirectHandler() http.HandlerFunc {
	repositoryRedirect := redirectSqlite.New(a.db)

	//cachedRepositoryRedirect = redirectSqlite.NewCached(repositoryRedirect, linkCache) todo

	serviceRedirect := redirectService.New(a.logger, repositoryRedirect, a.statManager)

	return redirectHandler.New(serviceRedirect, a.logger).ServeHTTP
}

func (a *App) setUpLogger(ctx context.Context) context.Context {
	var handler slog.Handler

	switch a.config.Env {
	case config.Local:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})
	case config.Dev:
		handler = slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	default:
		handler = slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true}))
	}

	handler = logs.NewContextMiddleware(handler)
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
	db, err := sqlutil.NewSQLiteStorage(a.config.DbConn)
	if err != nil {
		a.logger.ErrorContext(ctx, "cant connect to database", logs.ErrorMsg(err))
		os.Exit(1)
	}

	a.db = db

	a.logger.DebugContext(ctx, "database connected")

	err = sqlutil.Prepare(a.db)
	if err != nil {
		a.logger.ErrorContext(ctx, "cant prepare database", logs.ErrorMsg(err))
		os.Exit(1)
	}
	a.logger.DebugContext(ctx, "database prepared")
}

func (a *App) setUpStatisticManager() {
	repositoryStatistic := statistic.NewRepository(a.db)
	a.statManager = statistic.NewManager(1*time.Minute, 1000, repositoryStatistic, a.logger)
}
