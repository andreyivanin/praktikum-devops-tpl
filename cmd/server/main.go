package main

import (
	"devops-tpl/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	cfg := server.GetConfig()
	storage := server.InitConfig(cfg)
	r := server.NewRouter(storage)
	http.ListenAndServe(cfg.Address, r)

	sig := <-termSignal
	log.Panicln("Finished, reason:", sig.String())
}
