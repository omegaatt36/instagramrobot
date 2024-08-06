package main

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/omegaatt36/instagramrobot/app"
	"github.com/omegaatt36/instagramrobot/app/web"
	"github.com/omegaatt36/instagramrobot/health"
	"github.com/omegaatt36/instagramrobot/logging"
)

// Main is the entry point of the application.
func Main(ctx context.Context) {
	logging.Init(false)

	go health.StartServer()

	stopped := web.NewServer().Start(ctx)

	<-stopped
	<-ctx.Done()
	logging.Info("Shutting down")
}

func main() {
	app := app.App{
		Main:  Main,
		Flags: []cli.Flag{},
	}

	app.Run()
}
