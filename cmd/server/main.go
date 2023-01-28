package main

import (
	"devops-tpl/internal/server"
	"net/http"
)

func main() {
	cfg := server.GetEnvConfig()
	r := server.NewRouter()
	http.ListenAndServe(cfg.Address, r)
}
