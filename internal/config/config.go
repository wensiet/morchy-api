package config

import "github.com/caarlos0/env/v11"

var (
	AppName = "master"
	AppTag  = "v0.0.0"
)

type Config struct {
	Application struct {
		Port string `env:"APP_PORT,required"`
	}
	Database struct {
		Host     string `env:"DB_HOST,required"`
		User     string `env:"DB_USER,required"`
		Password string `env:"DB_PASSWORD,required"`
		Name     string `env:"DB_NAME,required"`
		Port     string `env:"DB_PORT,required"`
	}
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
