package handler

import (
	"devops-tpl/internal/storage/memstorage"
	"net/http"
	"testing"
)

func TestMetricJSON(t *testing.T) {

	// storage := memstorage.New()

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
			MetricJSON(tt.args.w, tt.args.r, tt.args.s)
		})
	}
}
