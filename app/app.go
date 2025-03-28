package app

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/omegaatt36/instagramrobot/cliflag"
)

// App provides a wrapper around urfave/cli app setup, including
// signal handling for graceful shutdown, panic recovery, and flag registration hooks.
type App struct {
	// Flags holds common CLI flags for the application.
	Flags []cli.Flag
	// Main is the core function of the application, executed after setup.
	Main func(ctx context.Context)
}

// before is executed by cli before the main Action. It initializes registered
// cliflag packages and includes basic panic recovery for the initialization phase.
func (a *App) before(c *cli.Context) (err error) {
	// Panic handling.
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered: ", r)
			debug.PrintStack()
			if err := cli.ShowAppHelp(c); err != nil {
				log.Fatal(err)
			}
			err = errors.New("init failed")
		}
	}()

	return cliflag.Initialize(c)
}

// after is executed by cli after the main Action. It finalizes registered
// cliflag packages (usually for cleanup).
func (a *App) after(c *cli.Context) error {
	return cliflag.Finalize(c)
}

// wrapMain wraps the execution of the application's Main function.
// It sets up a context cancellable by OS signals (SIGINT, SIGTERM)
// and includes panic recovery for the main application logic.
func (a *App) wrapMain(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("\nReceives signal: %v\n", sig)
		cancel()
	}()

	// Panic handling.
	defer func() {
		if r := recover(); r != nil {
			log.Println("Main recovered: ", r)
			debug.PrintStack()
		}
	}()

	a.Main(ctx)
	log.Println("terminated")

	return nil
}

// Run sets up the urfave/cli App with the configured flags, before/after hooks,
// and the main action wrapper, then executes it. It fatally logs any error
// during the cli app execution.
func (a *App) Run() {

	app := cli.NewApp()
	app.Flags = a.Flags
	app.Flags = append(app.Flags, cliflag.Globals()...)
	app.Before = a.before
	app.After = a.after
	app.Action = a.wrapMain

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
