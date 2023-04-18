package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/context"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Monitor) SendMetricsJSON() {
	url := CreateURLJSON(m.cfg)
	client := http.Client{}

	for _, metric := range m.Metrics {
		jsonMetric := Metrics{}

		switch metric.mtype {
		case "gauge":
			jsonMetric = Metrics{
				ID:    metric.name,
				MType: metric.mtype,
				Value: (*float64)(&metric.value),
			}
		case "counter":
			jsonMetric = Metrics{
				ID:    metric.name,
				MType: metric.mtype,
				Delta: (*int64)(&metric.delta),
			}

		}
		bytesMetric, err := json.Marshal(jsonMetric)
		if err != nil {
			panic(err)
		}

		body := bytes.NewBuffer(bytesMetric)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
		if err != nil {
			log.Fatalln(err)
		}

		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		response, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
		}

		// requestDump, err := httputil.DumpRequest(request, true)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		// fmt.Println(string(requestDump))

		if response != nil {
			fmt.Println("Status code", response.Status)

			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("Response body:\n %v\n", string(body))
		}

	}
}

func CreateURLJSON(cfg Config) string {
	var u url.URL
	u.Scheme = PROTOCOL
	u.Host = cfg.Address
	url := u.JoinPath("update")
	return url.String()
}
