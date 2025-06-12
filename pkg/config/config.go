package config

import (
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	MongoURI      string        `env:"MONGO_URI" envDefault:"mongodb://localhost:27017"`
	MongoDatabase string        `env:"MONGO_DATABASE" envDefault:"newsapp"`
	MongoTimeout  time.Duration `env:"MONGO_TIMEOUT" envDefault:"5s"`
	Server        struct {
		Address      string        `env:"SERVER_ADDRESS" envDefault:":8080"`
		ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"10s"`
		WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"10s"`
		IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"60s"`
	}
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
