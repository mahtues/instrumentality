package main

import (
	"fmt"
	"log"

	"net/http"

	"go.elastic.co/apm/module/apmhttp"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", HelloHandler())

	mux.Handle("/signin", SignInHandler())
	mux.Handle("/signup", SignUpHandler())
	mux.Handle("/signout", SignOutHandler())

	mux.Handle("/home", MustAuth(HomeHandler()))

	mux.Handle("/ping", PingHandler())
	port := 80

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), apmhttp.Wrap(mux)); err != nil {
		log.Fatal(err)
	}
}

func HelloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from app\n")
	})
}

func HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello %s, from app\n", r.Context().Value(nameKey).(string))
	})
}
