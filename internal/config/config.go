package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Telemetry TelemetryConfig `mapstructure:"telemetry"`
	Logger    LoggerConfig    `mapstructure:"logger"`
}

type ServerConfig struct {
	Host    string        `mapstructure:"host"`
	Port    int           `mapstructure:"port"`
	Timeout TimeoutConfig `mapstructure:"timeout"`
}

type TimeoutConfig struct {
	Read  time.Duration `mapstructure:"read"`
	Write time.Duration `mapstructure:"write"`
	Idle  time.Duration `mapstructure:"idle"`
}

type DatabaseConfig struct {
	Host     string     `mapstructure:"host"`
	Port     int        `mapstructure:"port"`
	User     string     `mapstructure:"user"`
	Password string     `mapstructure:"password"`
	Name     string     `mapstructure:"name"`
	SSLMode  string     `mapstructure:"sslmode"`
	Pool     PoolConfig `mapstructure:"pool"`
}

type PoolConfig struct {
	MaxOpen     int           `mapstructure:"max_open"`
	MaxIdle     int           `mapstructure:"max_idle"`
	MaxLifetime time.Duration `mapstructure:"max_lifetime"`
}

type LoggerConfig struct {
	Level       string `mapstructure:"level"`
	Environment string `mapstructure:"environment"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
