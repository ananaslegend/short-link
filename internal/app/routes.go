package app

import (
	"github.com/ananaslegend/short-link/internal/middleware"
	"github.com/ananaslegend/short-link/internal/statistic"
	"github.com/ananaslegend/short-link/pkg/clog"
	"github.com/go-pkgz/routegroup"
)

func (a *App) setUpRouter(router *routegroup.Bundle) {
	router.Use(clog.WithCtxLogger(a.logger))
	router.Use(middleware.WithRequestID)
	router.Use(middleware.WithRecover)
	router.Use(middleware.WithLoggingRequest)
	router.Use(statistic.WithStatisticRow(a.statManager))
	router.Use(statistic.WithSendingStatistic(a.statManager))

	router.HandleFunc("GET /{alias}", a.redirectHandler())
	router.HandleFunc("POST /link", a.saveLinkHandler())
}
