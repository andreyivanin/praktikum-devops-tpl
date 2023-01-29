package main

import (
	"devops-tpl/internal/server"
	"net/http"
)

func main() {
	server.InitFeatures()
	r := server.NewRouter()
	http.ListenAndServe(server.GetEnvConfig().Address, r)
}
