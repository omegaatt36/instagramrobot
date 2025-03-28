//go:generate go-enum
package logging

import (
	"github.com/urfave/cli/v2"
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
		EnvVars:     []string{"LOG_LEVEL"},
		Destination: &cfg.logLevel,
		Required:    false,
		DefaultText: zap.DebugLevel.String(),
		Value:       zap.DebugLevel.String(),
	})

	return flags
}
