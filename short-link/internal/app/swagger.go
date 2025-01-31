package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/ananaslegend/short-link/internal/config"
)

func (a *App) setupSwaggerDocumentationServer() {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", a.config.Swagger.Port))),
	)

	a.swaggerServer = &http.Server{
		Addr:    fmt.Sprintf(":%v", a.config.Swagger.Port),
		Handler: r,

		ReadHeaderTimeout: config.DefaultReadHeaderRequestTimeout,
		ReadTimeout:       config.DefaultReadRequestTimeout,
		WriteTimeout:      config.DefaultWriteTimeout,
		IdleTimeout:       config.DefaultIdleTimeout,
	}
}

func (a *App) runSwaggerDocumentationServer(ctx context.Context) error {
	if a.config.Env == config.Prod {
		return nil
	}

	a.logger.InfoContext(ctx, "swagger documentation server is started",
		slog.String("url", fmt.Sprintf("http://localhost:%v/swagger/index.html", a.config.Swagger.Port)),
	)

	if err := a.swaggerServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed start swagger documentation HTTP server, error: %w", err)
	}

	return nil
}
