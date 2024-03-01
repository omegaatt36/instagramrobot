//go:generate go-enum
package logging

import (
	"github.com/omegaatt36/instagramrobot/cliflag"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type cfg struct {
	logLevel string
}

var defaultConfig cfg

func init() {
	cliflag.Register(&defaultConfig)
}

// CliFlags returns cli flags to setup cache package.
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
