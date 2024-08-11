package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config is the object with configurable parameters
type Config struct {
	PostgresMigrate  bool   `env:"POSTGRES_MIGRATE" envDefault:"true"`
	PostgresHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort     int    `env:"POSTGRES_PORT,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDb       string `env:"POSTGRES_DB,required"`
	PostgresSslMode  string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`

	GrpcPort int `env:"GRPC_PORT" envDefault:"50051"`
}

// NewConfig is the function that returns Config object
func NewConfig() (*Config, error) {
	if err := loadEnv(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := cfg.readFromEnvironment(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) readFromEnvironment() error {
	return env.Parse(c)
}

func loadEnv() error {
	return godotenv.Load(".env")
}
