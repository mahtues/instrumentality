package main

import (
	"fmt"
	"log"

	"net/http"

	"go.elastic.co/apm/module/apmhttp"

	"github.com/kingnido/instrumentality/app/pkg/account"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", HelloHandler())
	mux.Handle("/signup", account.SignUpHandler())
	mux.Handle("/signin", account.SignInHandler())

	if err := http.ListenAndServe(fmt.Sprintf(":%d", 80), apmhttp.Wrap(mux)); err != nil {
		log.Fatal(err)
	}
}

func HelloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from app\n")
	})
}
