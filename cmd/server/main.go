package main

import (
	"devops-tpl/internal/server"
	"net/http"
)

func main() {
	cfg := server.GetEnvConfig()
	storage := server.InitConfig(cfg)
	r := server.NewRouter(storage)
	http.ListenAndServe(cfg.Address, r)
}
