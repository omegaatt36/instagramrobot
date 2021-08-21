package config

import (
	"github.com/feelthecode/instagramrobot/src/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ENV string

const (
	production  ENV = "prod"
	development ENV = "dev"
)

var (
	C *Config
)

type Config struct {
	APP_ENV   ENV    `mapstructure:"APP_ENV" validate:"required,oneof=prod dev"`
	BOT_TOKEN string `mapstructure:"BOT_TOKEN" validate:"required"`
}

func Load() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("$HOME/")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Panicf("could not unmarshal config file: %v", err)
	}

	// validate configuration keys and values
	if err := utils.Validate(C); err != nil {
		log.Panicf("config validation failed:\n%v", err)
	}

	log.WithField("APP_ENV", viper.GetString("APP_ENV")).Info("config loaded")
}

func IsDevelopment() bool {
	return C.APP_ENV == development
}

func IsProduction() bool {
	return C.APP_ENV == production
}
