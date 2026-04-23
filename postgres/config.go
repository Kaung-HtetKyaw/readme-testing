package postgres

import "github.com/caarlos0/env/v11"

var postgresConfigPrefix = "POSTGRES_"

type Config struct {
	User     string `env:"USER"`
	Password string `env:"PASSWORD,unset"`
	Host     string `env:"HOST" envDefault:"postgres"`
	Port     string `env:"PORT" envDefault:"5432"`
	DBName   string `env:"DB"`
	SSLMode  string `env:"SSLMODE"`
}

func NewConfigFromENV() (*Config, error) {
	opts := env.Options{
		Prefix:          postgresConfigPrefix,
		RequiredIfNoDef: true,
	}
	var c Config
	if err := env.ParseWithOptions(&c, opts); err != nil {
		return nil, err
	}
	return &c, nil
}
