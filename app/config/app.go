package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	envHttpPort = "CONFIG_HTTP_PORT"
)

type AppConfig struct {
	HttpPort int
}

func CreateAppConfig() (*AppConfig, error) {
	httpPort := viper.GetInt(envHttpPort)
	if httpPort == 0 {
		return nil, fmt.Errorf("missing required environment variable: %s\n", envHttpPort)
	}

	return &AppConfig{HttpPort: httpPort}, nil
}
