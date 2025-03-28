//go:generate go-enum -f=$GOFILE

package config

import (
	"github.com/urfave/cli/v3"

	"github.com/omegaatt36/instagramrobot/cliflag"
)

// Env defines the possible application environments (e.g., local, production).
// ENUM(
// local
// development
// production
// )
type Env string

// config holds the application configuration values, typically populated from flags or env vars.
type config struct {
	// appEnvironment stores the current application environment (e.g., "local", "production").
	appEnvironment string
}

var defaultConfig config

func init() {
	cliflag.Register(&defaultConfig)
}

// CliFlags returns the command-line flags associated with this configuration structure.
// These flags allow setting configuration values via command-line arguments or environment variables.
func (cfg *config) CliFlags() []cli.Flag {
	var flags []cli.Flag

	flags = append(flags, &cli.StringFlag{
		Name:        "app-env",
		Sources:     cli.EnvVars("APP_ENV"),
		Destination: &cfg.appEnvironment,
		Required:    false,
		DefaultText: EnvLocal.String(),
		Value:       EnvLocal.String(),
	})

	return flags
}

// IsLocal will return true if the APP_ENV is equals to local.
func IsLocal() bool {
	return defaultConfig.appEnvironment == EnvLocal.String()
}

// IsDevelopment will return true if the APP_ENV is equals to dev.
func IsDevelopment() bool {
	return defaultConfig.appEnvironment == EnvDevelopment.String()
}

// IsProduction will return true if the APP_ENV is equals to prod.
func IsProduction() bool {
	return defaultConfig.appEnvironment == EnvProduction.String()
}

// GetAppEnv returns the current app environment.
func GetAppEnv() string {
	return defaultConfig.appEnvironment
}
