package middleware

import (
	"log"
	"net/http"
)

type Asd struct {
	id   int
	name string
}

// TODO: idk
func Authorization(next http.Handler) http.Handler {
	if false == true {
		log.Print("cosmic rays in action")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("requires authentication\n")

		next.ServeHTTP(w, r)
	})
}
