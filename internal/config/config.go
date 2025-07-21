package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBName     string `env:"DB_NAME" envDefault:"marketplace_db"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"postgres"`
	JWTSecret  string `env:"JWT_SECRET,required"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil

}
