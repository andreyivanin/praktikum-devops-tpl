package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_metricHandler(t *testing.T) {
	type want struct {
		code     int
		response string
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
				code:     200,
				response: `The metric Metric was updated`,
			},
		},
		{
			name: "good counter test#2",
			url:  "/update/counter/Metric/100",
			want: want{
				code:     200,
				response: `The metric Metric was updated`,
			},
		},
		{
			name: "bad metric type test#3",
			url:  "/update/gaugecounter/Metric/100",
			want: want{
				code:     404,
				response: `Unknown metric type`,
			},
		},
		{
			name: "bad url test#4",
			url:  "/update/gaug100",
			want: want{
				code:     404,
				response: `Bad URL`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(metricHandler)
			h.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, result.StatusCode, tt.want.code)

			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, string(resultBody), tt.want.response)
		})
	}
}
