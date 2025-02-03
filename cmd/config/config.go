package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const ServiceName = "media-scout"

type General struct {
	ServiceName    string `mapstructure:"SERVICE_NAME"`
	TlsEnabled     bool   `mapstructure:"TLS_ENABLED"`
	LoggingEnabled bool   `mapstructure:"LOGGING_ENABLED"`
}
type Config struct {
	General General `mapstructure:"GENERAL,squash"`
	HTTP    HTTP    `mapstructure:"HTTP"`
	DB      DB      `mapstructure:"DB"`
}

type HTTP struct {
	Port             int           `mapstructure:"PORT"`
	GracefulShutdown time.Duration `mapstructure:"GRACEFUL_SHUTDOWN"`
}

type DB struct {
	DSN                    string        `mapstructure:"DSN"`
	MaxOpenConnections     int           `mapstructure:"MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections     int           `mapstructure:"MAX_IDLE_CONNECTIONS"`
	MaxConnectionsLifetime time.Duration `mapstructure:"MAX_CONNECTIONS_LIFETIME"`
}

func NewConfig(path, filename string) (Config, error) {
	vpr := viper.NewWithOptions(viper.KeyDelimiter("__"))
	vpr.AddConfigPath(path)
	vpr.SetConfigName(filename)
	vpr.SetConfigType("env")
	vpr.AutomaticEnv()
	if err := vpr.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := vpr.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}
