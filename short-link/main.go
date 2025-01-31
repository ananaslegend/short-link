package main

import (
	"context"
	"os/signal"
	"syscall"

	_ "github.com/ananaslegend/short-link/docs"
	"github.com/ananaslegend/short-link/internal/app"
)

//	@title		Short link service
//	@version	v1.0

// @host		localhost:8080
// @BasePath	/api/v1
// @schemes	http https
func main() {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a := app.New(ctx)

	go func() {
		if err := a.Run(ctx); err != nil {
			cancel()
		}
	}()

	<-ctx.Done()

	a.Close()
}
