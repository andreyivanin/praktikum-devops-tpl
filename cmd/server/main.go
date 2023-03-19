package main

import (
	"context"
	"log"
	"net/http"

	"devops-tpl/internal/server"
)

func main() {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.InitSignal(ctx)

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
