package config

import (
	"flag"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ENV string

const (
	production  ENV = "prod"
	development ENV = "dev"
)

type config struct {
	APP_ENV   ENV    `mapstructure:"APP_ENV" validate:"required,oneof=prod dev"`
	BOT_TOKEN string `mapstructure:"BOT_TOKEN" validate:"required"`
}

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
		log.Info("Reading config from environment variables")
		viper.BindEnv("APP_ENV")
		viper.BindEnv("BOT_TOKEN")
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

func IsDevelopment() bool {
	return ENV(viper.GetString("APP_ENV")) == development
}

func IsProduction() bool {
	return ENV(viper.GetString("APP_ENV")) == production
}
