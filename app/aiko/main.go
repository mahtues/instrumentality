package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mahtues/instrumentality/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	hostname := os.Getenv("HOSTNAME")
	log.Warningf("hello world")

	http.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "pong from %s\n", hostname)
	})

	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":2112", nil)
}
