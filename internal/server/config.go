package server

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	SERVERADDRPORT = "localhost:8080"
	STOREINTERVAL  = 0
	STOREFILE      = "devops-metrics-db.json"
	RESTORE        = true
)

type Config struct {
	Address          string        `env:"ADDRESS"`
	StoreInterval    time.Duration `env:"STORE_INTERVAL"`
	StoreFile        string        `env:"STORE_FILE"`
	RestoreSavedData bool          `env:"RESTORE"`
}

func GetEnvConfig() Config {
	var cfg = Config{
		Address:          SERVERADDRPORT,
		StoreInterval:    STOREINTERVAL * time.Second,
		StoreFile:        STOREFILE,
		RestoreSavedData: RESTORE,
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg

}
