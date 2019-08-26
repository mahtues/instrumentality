package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"net/http"
)

func main() {
	mux2 := http.NewServeMux()
	mux2.Handle("/a", LookupEnvMiddleware("TEST", PrintHandler("a")))
	mux2.Handle("/b", PrintHandler("b"))
	mux2.Handle("/", PrintHandler("other"))

	mux := http.NewServeMux()
	mux.Handle("/", DurationMiddleware(EchoHandler()))
	mux.Handle("/private", AuthorizationMiddleware(EchoHandler()))
	mux.Handle("/public/", http.StripPrefix("/public", mux2))

	err := http.ListenAndServe(":80", mux)

	if err != nil {
		log.Fatal(err)
	}
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("requires authentication")

		next.ServeHTTP(w, r)
	})
}

func DurationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Println("elapsed:", time.Since(start))
	})
}

func LookupEnvMiddleware(key string, next http.Handler) http.Handler {
	value := os.Getenv(key)

	if value == "" {
		log.Printf("env var '%s' not found or empty. this middleware will be dropped", key)
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("key = %s, name = %s", key, value)
		next.ServeHTTP(w, r)
	})
}

func EchoHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, r.Body)
	})
}

func PrintHandler(text string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, text)
	})
}
