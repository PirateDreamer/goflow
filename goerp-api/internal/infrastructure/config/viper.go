package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Email    EmailConfig
	Swagger  SwaggerConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type EmailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type SwaggerConfig struct {
	User     string
	Password string
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	DSN string
}

func InitConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	return &cfg, nil
}
