package main

import (
	"context"
	"flag"
	"github.com/ananaslegend/short-link/internal/app"
	"os/signal"
	"syscall"
)

func main() {
	confPath := flag.String("config", "../../config/app-config.yml", "path to config file")
	flag.Parse()

	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a := app.New(ctx, *confPath)

	go func() {
		if err := a.Run(); err != nil {
			cancel()
		}
	}()

	<-ctx.Done()

	a.Close()
}
