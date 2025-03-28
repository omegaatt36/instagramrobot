//go:generate go-enum
package logging

import (
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"

	"github.com/omegaatt36/instagramrobot/cliflag"
)

// cfg holds logging configuration, primarily the log level.
type cfg struct {
	// logLevel specifies the minimum level for log messages (e.g., "debug", "info", "error").
	logLevel string
}

var defaultConfig cfg

func init() {
	cliflag.Register(&defaultConfig)
}

// CliFlags returns the command-line flags for configuring the logging package.
func (cfg *cfg) CliFlags() []cli.Flag {
	var flags []cli.Flag

	flags = append(flags, &cli.StringFlag{
		Name:        "log-level",
		Sources:     cli.EnvVars("LOG_LEVEL"),
		Destination: &cfg.logLevel,
		Required:    false,                   // Required is still supported
		DefaultText: zap.DebugLevel.String(), // DefaultText is still supported
		Value:       zap.DebugLevel.String(), // Value sets the default
	})

	return flags
}
