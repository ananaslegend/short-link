package app

import (
	"github.com/go-pkgz/routegroup"

	_ "github.com/ananaslegend/short-link/docs"
	"github.com/ananaslegend/short-link/internal/middleware"
	"github.com/ananaslegend/short-link/internal/statistic"
	"github.com/ananaslegend/short-link/pkg/cslog"
)

func (a *App) setUpRouter(router *routegroup.Bundle) {
	router.Use(cslog.Middleware(a.logger))
	router.Use(middleware.WithRequestID)
	router.Use(middleware.WithRecover)
	router.Use(middleware.WithLoggingRequest)
	router.Use(statistic.WithStatisticRow(a.statManager))
	router.Use(statistic.WithSendingStatistic(a.statManager))

	router.HandleFunc("GET /api/v1/{alias}", a.redirectHandler())
	router.HandleFunc("POST /api/v1/link", a.saveLinkHandler())
}
