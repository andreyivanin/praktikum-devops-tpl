package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"devops-tpl/internal/server"
)

func InitSignal(ctx context.Context) {
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-termSignal
	log.Println("Finished, reason:", sig.String())
	os.Exit(0)
}

func main() {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go InitSignal(ctx)

	cfg, err := server.GetConfig()
	if err != nil {
		log.Println(err)
	}

	storage, err := server.InitStorage(cfg)
	if err != nil {
		log.Println(err)
	}

	r, err := server.NewRouter(storage)
	if err != nil {
		log.Println(err)
	}

	err = http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Println(err)
	}
}
