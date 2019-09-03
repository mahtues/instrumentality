package main

import (
	"log"

	"go.elastic.co/apm/module/apmhttp"
	"net/http"

	"github.com/kingnido/instrumentality/handler"
	"github.com/kingnido/instrumentality/middleware"
)

func main() {
	mux2 := http.NewServeMux()
	mux2.Handle("/b", handler.Print("b"))
	mux2.Handle("/", handler.Print("other"))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Duration(handler.Echo()))
	mux.Handle("/private", middleware.Authorization(handler.Echo()))
	mux.Handle("/public/", http.StripPrefix("/public", mux2))

	mux.Handle("/shortest_path", handler.ShortestPath())
	mux.Handle("/shortest_path_v2", handler.ShortestPathV2())

	err := http.ListenAndServe(":80", apmhttp.Wrap(mux))

	if err != nil {
		log.Fatal(err)
	}
}
