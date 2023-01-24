package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewRouter() chi.Router {
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
		r.Route("/{mtype}/{mname}", func(r chi.Router) {
			r.Get("/", MetricGetHandler)
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", MetricSummaryHandler)
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

func Test_MetricUpdateHandler(t *testing.T) {
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
			r := NewRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()

			code, body := testRequest(t, ts, "POST", tt.url)
			assert.Equal(t, tt.want.code, code)
			assert.Equal(t, tt.want.body, body)

			// request := httptest.NewRequest(http.MethodPost, tt.url, nil)
			// w := httptest.NewRecorder()
			// h := http.HandlerFunc(metricUpdateHandler)
			// h.ServeHTTP(w, request)
			// result := w.Result()
			// assert.Equal(t, result.StatusCode, tt.want.code)

			// defer result.Body.Close()
			// resultBody, err := io.ReadAll(result.Body)
			// if err != nil {
			// 	t.Fatal(err)
			// }

			// assert.Equal(t, string(resultBody), tt.want.response)
		})
	}
}
