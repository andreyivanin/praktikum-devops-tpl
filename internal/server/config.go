package server

import (
	"log"
	"os"
	"time"

	"devops-tpl/internal/storage/file"
	"devops-tpl/internal/storage/memory"

	"github.com/caarlos0/env/v6"
)

const (
	SERVERADDRPORT = "localhost:8080"
	STOREINTERVAL  = 10
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

func InitConfig() {
	cfg := GetEnvConfig()
	if cfg.RestoreSavedData {
		reader, err := file.NewReader(cfg.StoreFile)
		if err != nil {
			log.Fatal(err)
		}

		checkFile, err := os.Stat(cfg.StoreFile)
		if err != nil {
			log.Fatal(err)
		}

		size := checkFile.Size()

		if size == 0 {
			writer, err := file.NewWriter(cfg.StoreFile)
			if err != nil {
				log.Fatal(err)
			}

			if err := writer.WriteDatabase(); err != nil {
				log.Fatal(err)
			}
		}

		if file.DB, err = reader.ReadDatabase(); err != nil {
			log.Fatal(err)
		}
	}

	if cfg.StoreFile != " " {
		if cfg.StoreInterval == 0 {
			for range memory.MetricUpdated {
				writer, err := file.NewWriter(cfg.StoreFile)
				if err != nil {
					log.Fatal(err)
				}

				if err := writer.WriteDatabase(); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			ticker := time.NewTicker(cfg.StoreInterval)
			for range ticker.C {
				writer, err := file.NewWriter(cfg.StoreFile)
				if err != nil {
					log.Fatal(err)
				}

				if err := writer.WriteDatabase(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
