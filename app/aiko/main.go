package main

import (
	"fmt"
	"log"
	"net/http"

	"go.elastic.co/apm/module/apmhttp"

	"github.com/mahtues/instrumentality/server"
)

func main() {
	config := server.Config{
		Port: 8080,
	}

	srv, err := server.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), apmhttp.Wrap(srv)); err != nil {
		log.Fatal(err)
	}
}
