package main

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/omegaatt36/instagramrobot/app"
	"github.com/omegaatt36/instagramrobot/app/bot/config"
	"github.com/omegaatt36/instagramrobot/app/web"
	"github.com/omegaatt36/instagramrobot/health"
	"github.com/omegaatt36/instagramrobot/logging"
)

// Main is the entry point of the application.
func Main(ctx context.Context, cmd *cli.Command) error {
	logging.Init(!config.IsLocal())

	go health.StartServer()

	// Wait for the web server to stop and for the context cancellation signal.
	stopped := web.NewServer().Start(ctx)

	select {
	case <-stopped:
		logging.Info("Web server stopped normally.")
	case <-ctx.Done():
		logging.Info("Shutdown signal received, waiting for web server to stop...")
		// Wait for server to finish stopping after context cancellation
		<-stopped
		logging.Info("Web server stopped after signal.")
	}

	logging.Info("Shutting down")
	return nil // Indicate successful execution
}

func main() {
	// Initialize app structure for v3
	application := app.App{
		Main: Main,
		// Base flags for the app itself, if any, go here.
		// Module flags are added via cliflag.Globals() in app.Run()
		Flags: []cli.Flag{},
	}

	application.Run()
}
