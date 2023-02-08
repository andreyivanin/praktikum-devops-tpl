package main

import (
	"devops-tpl/internal/server"
	"net/http"
)

func main() {
	server.InitConfig()
	r := server.NewRouter()
	http.ListenAndServe(server.GetEnvConfig().Address, r)
}
