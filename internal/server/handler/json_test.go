package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"devops-tpl/internal/storage/memstorage"
)

func TestMetricJSON(t *testing.T) {

	type args struct {
		w http.ResponseWriter
		r *http.Request
		s *memstorage.MemStorage
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := memstorage.New()
			r := NewRouter(storage)
			ts := httptest.NewServer(r)
			defer ts.Close()
		})
	}
}
