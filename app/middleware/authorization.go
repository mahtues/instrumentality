package middleware

import (
	"log"
	"net/http"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("requires authentication")

		next.ServeHTTP(w, r)
	})
}
