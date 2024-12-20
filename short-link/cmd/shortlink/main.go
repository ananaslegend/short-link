package main

import (
	"context"
	"github.com/ananaslegend/short-link/internal/app"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a := app.New(ctx)

	go func() {
		if err := a.Run(); err != nil {
			cancel()
		}
	}()

	<-ctx.Done()

	a.Close()
}
