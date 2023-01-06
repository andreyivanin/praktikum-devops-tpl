package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", metricHandler)
	http.ListenAndServe(":8080", nil)
}
