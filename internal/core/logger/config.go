package core_logger

//Настройка логгера

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `envconfig:"LEVEL" default:"DEBUG"`
	Folder string `envconfig:"FOLDER" default:"out/logs"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process logger config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get logger config: %w", err))
	}

	return config
}
