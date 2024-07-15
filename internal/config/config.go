package config

import (
	"log"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type Config struct {
	Env        Env              `env:"ENV" env-required:"true"`
	ClickHouse DBConfig         `env-prefix:"CLICKHOUSE_"`
	Postgres   DBConfig         `env-prefix:"POSTGRES_"`
	HTTPServer HTTPServerConfig `env-prefix:"HTTP_SERVER_"`
}

type HTTPServerConfig struct {
	IpAddress string        `env:"IP_ADDRESS" env-required:"true"`
	Port      string        `env:"PORT" env-required:"true"`
	Timeout   time.Duration `env:"TIMEOUT" env-default:"4s"`
}

type DBConfig struct {
	Host     string `env:"HOST" env-default:"localhost"`
	Port     string `env:"PORT" env-required:"true"`
	DBName   string `env:"DB_NAME" env-required:"true"`
	Username string `env:"USERNAME" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	SSLMode  string `env:"SSL_MODE"`
}

var (
	once sync.Once
	cfg  Config
)

func MustNew() Config {
	once.Do(func() {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("could not read config: %s", err)
		}
	})

	return cfg
}
