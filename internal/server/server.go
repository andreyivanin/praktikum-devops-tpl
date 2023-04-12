package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"devops-tpl/internal/server/handler"
	"devops-tpl/internal/storage"
	"devops-tpl/internal/storage/filestorage"
	"devops-tpl/internal/storage/memstorage"
)

func NewRouter(storage storage.Storage) (chi.Router, error) {
	handler := handler.NewHandler(storage)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handler.MetricJSON)
		r.Route("/{mtype}/{mname}/{mvalue}", func(r chi.Router) {
			r.Post("/", handler.MetricUpdate)
			r.Get("/", handler.MetricUpdate)
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", handler.MetricSummaryJSON)
		r.Route("/{mtype}/{mname}", func(r chi.Router) {
			r.Get("/", handler.MetricGet)
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.MetricSummary)
	})

	return r, nil
}

func newMemoryStorage() *memstorage.MemStorage {
	storage := memstorage.New()
	return storage
}

func newFileStorage(cfg Config) *filestorage.FileStorage {
	storage := filestorage.New(cfg.StoreFile)
	if cfg.StoreInterval != 0 {
		// ctx, cancel := context.WithCancel(context.Background())
		// defer cancel()
		// storage.SaveTicker(ctx, cfg.StoreInterval)
		go storage.SaveTicker(cfg.StoreInterval)

	} else {
		storage.SyncMode = true
	}

	if cfg.RestoreSavedData {
		storage.Restore(cfg.StoreFile)
	}
	return storage
}

func NewStorage(cfg Config) (storage.Storage, error) {
	if cfg.StoreFile != " " {
		return newFileStorage(cfg), nil
	} else {
		return newMemoryStorage(), nil
	}

}
