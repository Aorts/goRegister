package config

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
