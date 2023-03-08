package main

import (
	"devops-tpl/internal/server"
	"net/http"
)

func main() {
	go server.InitSignal()
	cfg := server.GetConfig()
	storage := server.InitConfig(cfg)
	r := server.NewRouter(storage)
	http.ListenAndServe(cfg.Address, r)
}
