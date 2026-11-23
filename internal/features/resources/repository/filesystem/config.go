package resources_filesystem_repository

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Folder string `envconfig:"FOLDER" default:"out/resource-files"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("RESOURCE_FILES", &config); err != nil {
		return Config{}, fmt.Errorf("process resource file storage config: %w", err)
	}

	return config, nil
}
