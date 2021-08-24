package config

import (
	"flag"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type env string

const (
	production  env = "prod"
	development env = "dev"
)

type config struct {
	APP_ENV   env    `mapstructure:"APP_ENV" validate:"required,oneof=prod dev"`
	BOT_TOKEN string `mapstructure:"BOT_TOKEN" validate:"required"`
}

// Load the configuration file
func Load() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	configPath := flag.String("config-path", "", "config file path")
	flag.Parse()

	if *configPath != "" {
		// Register config folder from flag
		viper.AddConfigPath(*configPath)
	} else {
		// Current directory
		viper.AddConfigPath("$HOME/")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		// log.Warn(err)
		log.Info("Config loaded from env vars")

		// TODO: refactor with foreach
		if err := viper.BindEnv("APP_ENV"); err != nil {
			log.Fatal("Couldn't bind the APP_ENV env var")
		}
		if err := viper.BindEnv("BOT_TOKEN"); err != nil {
			log.Fatal("Couldn't bind the BOT_TOKEN env var")
		}
	}

	var c *config

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("Could not unmarshal config: %v", err)
	}

	// validate configuration keys and values
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		log.Fatalf("Config validation failed:\n%v", err)
	}

	log.WithField("APP_ENV", viper.GetString("APP_ENV")).Info("Config loaded")
}

// IsDevelopment will return true if the APP_ENV is equals to dev
func IsDevelopment() bool {
	return env(viper.GetString("APP_ENV")) == development
}

// IsProduction will return true if the APP_ENV is equals to prod
func IsProduction() bool {
	return env(viper.GetString("APP_ENV")) == production
}
