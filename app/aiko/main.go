package main

import (
	"fmt"
	"net/http"

	"github.com/mahtues/instrumentality/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Warningf("hello world")

	http.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":2112", nil)
}
