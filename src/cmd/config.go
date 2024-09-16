package cmd

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS,required"`
	Database      DatabaseConfig
}

type DatabaseConfig struct {
	Conn     string `env:"POSTGRES_CONN,required"`
	JDBCURL  string `env:"POSTGRES_JDBC_URL,required"`
	User     string `env:"POSTGRES_USERNAME,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	Name     string `env:"POSTGRES_DATABASE,required"`
}

func Load() (*Config, error) {
	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
