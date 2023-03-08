package agent

import (
	"testing"
	"time"
)

func TestSendMetricsJSON(t *testing.T) {
	var cfg = Config{
		Address:        SERVERADDRPORT,
		PollInterval:   POLLINTERVAL * time.Second,
		ReportInterval: REPORTINTERVAL * time.Second,
	}

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendMetricsJSON(cfg)
		})
	}
}
