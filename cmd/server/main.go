package main

import (
	"devops-tpl/internal/server"
	"net/http"
)

func main() {
	r := server.NewRouter()
	http.ListenAndServe(":8080", r)
}
