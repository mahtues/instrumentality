package main

import (
	"fmt"
	"log"

	"net/http"

	"go.elastic.co/apm/module/apmhttp"

	"github.com/kingnido/instrumentality/api"
)

func main() {
	mux := http.NewServeMux()


	mux.Handle("/", helloHandler())
	mux.Handle("/signup", api.SignUpHandler())
	mux.Handle("/signin", api.SignInHandler())

	if err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), apmhttp.Wrap(mux)); err != nil {
		log.Fatal(err)
	}
}

func helloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from app\n")
	})
}
