package server

import (
	"devops-tpl/internal/server/handlers"
	"devops-tpl/internal/storage"
	"log"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handlers.MetricJSONHandler)
		r.Route("/{mtype}/{mname}/{mvalue}", func(r chi.Router) {
			r.Post("/", handlers.MetricUpdateHandler)
			r.Get("/", handlers.MetricUpdateHandler)
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", handlers.MetricSummaryJSONHandler)
		r.Route("/{mtype}/{mname}", func(r chi.Router) {
			r.Get("/", handlers.MetricGetHandler)
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.MetricSummaryHandler)
	})

	return r
}

func InitFeatures() {
	cfg := GetEnvConfig()
	if cfg.RestoreSavedData == true {
		reader, err := storage.NewReader(GetEnvConfig().StoreFile)
		if err != nil {
			log.Fatal(err)
		}

		if err := reader.ReadDatabase(); err != nil {
			log.Fatal(err)
		}
	}

	if cfg.StoreFile != " " {
		go StoreOnDisk(cfg)
	}
}

func StoreOnDisk(cfg Config) {
	if cfg.StoreInterval == 0 {
		for {
			select {
			case <-storage.MetricUpdated:
				writer, err := storage.NewWriter(GetEnvConfig().StoreFile)
				if err != nil {
					log.Fatal(err)
				}

				if err := writer.WriteDatabase(); err != nil {
					log.Fatal(err)
				}
			}
		}
	} else {
		ticker := time.NewTicker(cfg.StoreInterval)
		for {
			select {
			case <-ticker.C:
				writer, err := storage.NewWriter(GetEnvConfig().StoreFile)
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
