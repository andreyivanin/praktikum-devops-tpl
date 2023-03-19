package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"devops-tpl/internal/server/handler"
	"devops-tpl/internal/storage"
	"devops-tpl/internal/storage/filestorage"
	"devops-tpl/internal/storage/memstorage"
)

func NewRouter(db storage.Storage) (chi.Router, error) {
	MetricJSONHandler := func(w http.ResponseWriter, r *http.Request) {
		handler.MetricJSON(w, r, db)
	}

	MetricSummaryJSONHandler := func(w http.ResponseWriter, r *http.Request) {
		handler.MetricSummaryJSON(w, r, db)
	}

	MetricUpdateHandler := func(w http.ResponseWriter, r *http.Request) {
		handler.MetricUpdate(w, r, db)
	}

	MetricGetHandler := func(w http.ResponseWriter, r *http.Request) {
		handler.MetricGet(w, r, db)
	}

	MetricSummaryHandler := func(w http.ResponseWriter, r *http.Request) {
		handler.MetricSummary(w, r, db)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Route("/update", func(r chi.Router) {
		r.Post("/", MetricJSONHandler)
		r.Route("/{mtype}/{mname}/{mvalue}", func(r chi.Router) {
			r.Post("/", MetricUpdateHandler)
			r.Get("/", MetricUpdateHandler)
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", MetricSummaryJSONHandler)
		r.Route("/{mtype}/{mname}", func(r chi.Router) {
			r.Get("/", MetricGetHandler)
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", MetricSummaryHandler)
	})

	return r, nil
}

func runMemoryStorage() *memstorage.MemStorage {
	storage := memstorage.New()
	return storage
}

func runFileStorage(cfg Config) *filestorage.FileStorage {
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

func InitStorage(cfg Config) (storage.Storage, error) {
	if cfg.StoreFile != " " {
		return runFileStorage(cfg), nil
	} else {
		return runMemoryStorage(), nil
	}

}

func InitSignal(ctx context.Context) {
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-termSignal
	log.Println("Finished, reason:", sig.String())
	os.Exit(0)
}
