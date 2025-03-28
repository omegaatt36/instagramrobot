package cliflag

import (
	"log"

	"github.com/urfave/cli/v2"
)

// CliFlager is interface to describe a struct
// which is a set of options to singleton in package.
// This struct has method CliFlags to returns the options this package needed.
// CliFlager defines the interface for packages that provide command-line flags.
// This allows central registration and collection of flags from various modules.
type CliFlager interface {
	// CliFlags returns a slice of cli.Flag definitions for the package.
	CliFlags() []cli.Flag
}

// Beforer is interface for some package may needs an before hook to init
// or validates before calls main function.
// Beforer defines an optional interface for packages that need an initialization
// hook to run *before* the main application logic, after flags are parsed.
type Beforer interface {
	// Before performs initialization tasks. It receives the cli context.
	Before(*cli.Context) error
}

// Afterer is interface for some package may needs an after hook to destroy
// or cleanup after main function exited.
// Afterer defines an optional interface for packages that need a cleanup
// hook to run *after* the main application logic has finished.
type Afterer interface {
	// After performs cleanup tasks.
	After()
}

var cliFlagers []CliFlager

// Register adds a CliFlager instance to the global registry. This ensures its
// flags are collected and its Before/After hooks (if implemented) are called.
func Register(f CliFlager) {
	cliFlagers = append(cliFlagers, f)
}

// IsBeforer checks interface conversion.
func IsBeforer(Beforer) {}

// IsAfterer checks interface conversion.
func IsAfterer(Afterer) {}

// Globals collects and returns a combined slice of all cli.Flag definitions
// from all registered CliFlager instances.
func Globals() []cli.Flag {
	var res []cli.Flag
	for _, f := range cliFlagers {
		res = append(res, f.CliFlags()...)
	}

	return res
}

// Initialize iterates through all registered packages and calls the Before method
// for those that implement the Beforer interface.
func Initialize(c *cli.Context) error {
	for _, f := range cliFlagers {
		b, ok := f.(Beforer)
		if ok {
			log.Printf("running %T", b)
			err := b.Before(c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Finalize iterates through all registered packages and calls the After method
// (using defer for LIFO execution) for those that implement the Afterer interface.
func Finalize(c *cli.Context) error {
	//revive:disable:defer
	for _, f := range cliFlagers {
		a, ok := f.(Afterer)
		if ok {
			defer a.After()
		}
	}
	//revive:enable:defer

	return nil
}
