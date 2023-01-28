package server

import (
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	SERVERADDRPORT = "localhost:8080"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func GetEnvConfig() Config {
	var cfg = Config{
		Address: SERVERADDRPORT,
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg

}
