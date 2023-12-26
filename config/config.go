package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Version string
	Server  struct {
		Port string
	}
	Database DatabaseConfig
	Redis    RedisConfig
}

type DatabaseConfig struct {
	Driver   string
	Username string
	Password string
	Host     string
	Database string
}

type RedisConfig struct {
	Host string
}

func InitConfig() (*Config, error) {
	viper.SetDefault("Server.Port", ":8080")
	viper.SetDefault("Database.Version", "v0.0.1")
	viper.SetDefault("Database.Driver", "postgres")
	viper.SetDefault("Database.Username", "ts")
	viper.SetDefault("Database.Password", "ts")
	viper.SetDefault("Database.Host", "localhost:5432")
	viper.SetDefault("Database.Database", "postgres")
	viper.SetDefault("Redis.Host", "localhost:6379")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
