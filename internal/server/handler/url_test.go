package handler

import (
	"devops-tpl/internal/storage"
	"devops-tpl/internal/storage/memstorage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewRouter(db storage.Storage) chi.Router {
	handler := NewHandler(db)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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

	return r
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func Test_MetricUpdate(t *testing.T) {
	type want struct {
		code int
		body string
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "good gauge test#1",
			url:  "/update/gauge/Metric/100",
			want: want{
				code: 200,
				body: `The metric Metric was updated`,
			},
		},
		{
			name: "good counter test#2",
			url:  "/update/counter/Metric/100",
			want: want{
				code: 200,
				body: `The metric Metric was updated`,
			},
		},
		{
			name: "bad metric type test#3",
			url:  "/update/gaugecounter/Metric/100",
			want: want{
				code: 501,
				body: `Bad metric type`,
			},
		},
		{
			name: "bad url test#4",
			url:  "/update/gaug100",
			want: want{
				code: 404,
				body: `404 page not found
`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := memstorage.New()
			r := NewRouter(storage)
			ts := httptest.NewServer(r)
			defer ts.Close()

			code, body := testRequest(t, ts, "POST", tt.url)
			assert.Equal(t, tt.want.code, code)
			assert.Equal(t, tt.want.body, body)

		})
	}
}
