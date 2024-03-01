//go:generate go-enum -f=$GOFILE

package config

import (
	"github.com/omegaatt36/instagramrobot/app/cliflag"
	"github.com/urfave/cli/v2"
)

// Env represents the environment of the application.
// ENUM(
// local
// development
// production
// )
type Env string

type config struct {
	appEnvironment string
	botToken       string
}

var defaultConfig config

func init() {
	cliflag.Register(&defaultConfig)
}

// CliFlags returns cli flags to setup cache package.
func (cfg *config) CliFlags() []cli.Flag {
	var flags []cli.Flag

	flags = append(flags, &cli.StringFlag{
		Name:        "app-env",
		EnvVars:     []string{"APP_ENV"},
		Destination: &cfg.appEnvironment,
		Required:    false,
		DefaultText: EnvLocal.String(),
		Value:       EnvLocal.String(),
	})

	flags = append(flags, &cli.StringFlag{
		Name:        "bot-token",
		EnvVars:     []string{"BOT_TOKEN"},
		Destination: &cfg.botToken,
		Required:    true,
	})

	return flags
}

// BotToken returns the bot token.
func BotToken() string {
	return defaultConfig.botToken
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
