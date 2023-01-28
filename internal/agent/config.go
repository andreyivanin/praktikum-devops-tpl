package agent

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	PROTOCOL       = "http"
	SERVERADDRPORT = "localhost:8080"
	POLLINTERVAL   = 2
	REPORTINTERVAL = 10
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

func GetEnvConfig() Config {
	var cfg = Config{
		Address:        SERVERADDRPORT,
		PollInterval:   POLLINTERVAL * time.Second,
		ReportInterval: REPORTINTERVAL * time.Second,
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg

}
