package main

import (
	"net/http"
)

func main() {
	r := NewRouter()
	http.ListenAndServe(":8080", r)
}
