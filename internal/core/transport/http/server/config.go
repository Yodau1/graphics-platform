package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// файл читает адрес сервера и время,
// которое даётся активным запросам на завершение

type Config struct {
	// Addr — адрес и порт сервера, например ":5050" (слушать на всех интерфейсах).
	Addr string `envconfig:"ADDR" required:"true"`
	// ShutdownTimeout — время ожидания завершения активных запросов при остановке сервера.
	// По истечении этого времени сервер принудительно закрывает соединения.
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

// NewConfig читает конфигурацию сервера из переменных окружения.
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("HTTP", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get HTTP server config: %w", err))
	}
	return config
}
