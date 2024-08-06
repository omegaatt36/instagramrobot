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

// App is cli wrapper that do some common operation and creates signal handler.
type App struct {
	Flags []cli.Flag
	Main  func(ctx context.Context)
}

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

func (a *App) after(c *cli.Context) error {
	return cliflag.Finalize(c)
}

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

// Run setups everything and runs Main.
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
