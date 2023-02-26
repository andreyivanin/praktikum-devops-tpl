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
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func SendMetricsJSON() {
	url := CreateURLJSON()
	client := http.Client{}

	gMetrics := GMetricGeneratorNew()
	for _, gMetric := range gMetrics {
		name := gMetric.Name
		value := gMetric.Value
		jsonMetric := Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}
		metricJSON, err := json.Marshal(jsonMetric)
		if err != nil {
			panic(err)
		}

		body := bytes.NewBuffer(metricJSON)
		request, err := http.NewRequest(http.MethodPost, url, body)
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

	cMetrics := CMetricGeneratorNew()
	for _, cMetric := range cMetrics {
		name := cMetric.Name
		value := cMetric.Value
		jsonMetric := Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
		}

		metricJSON, err := json.Marshal(jsonMetric)
		if err != nil {
			panic(err)
		}

		body := bytes.NewBuffer(metricJSON)
		request, err := http.NewRequest(http.MethodPost, url, body)
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

func CreateURLJSON() string {
	var u url.URL
	u.Scheme = PROTOCOL
	u.Host = GetConfig().Address
	url := u.JoinPath("update")
	return url.String()
}
