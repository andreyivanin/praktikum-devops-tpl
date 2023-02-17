package server

import (
	"devops-tpl/internal/server/handler"
	"devops-tpl/internal/storage"
	"devops-tpl/internal/storage/filestorage"
	"devops-tpl/internal/storage/memstorage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db storage.Storage) chi.Router {
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

	return r
}

func RunMemory() *memstorage.MemStorage {
	storage := memstorage.New()
	return storage
}

func RunFile(cfg Config) *filestorage.FileStorage {
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

func InitConfig(cfg Config) storage.Storage {
	if cfg.StoreFile != " " {
		return RunFile(cfg)
	} else {
		return RunMemory()
	}

}
