package main

import (
	"fmt"
	"log"

	"net/http"

	"go.elastic.co/apm/module/apmhttp"

	"github.com/kingnido/instrumentality/pkg/account"
	"github.com/kingnido/instrumentality/pkg/api"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", HelloHandler())
	mux.Handle("/signup", api.SignUpHandler(account.DefaultCreater()))
	mux.Handle("/signin", api.SignInHandler(account.DefaultVerifier()))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", 80), apmhttp.Wrap(mux)); err != nil {
		log.Fatal(err)
	}
}

func HelloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from app\n")
	})
}
