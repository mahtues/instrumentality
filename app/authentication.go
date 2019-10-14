package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type contextKey int

const nameKey = contextKey(0)

var AccountStore = map[string]string{}
var CookieStore = map[string]string{}

func SignUpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if len(username) == 0 || len(password) == 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if _, ok := AccountStore[username]; ok {
			http.Error(w, "Account already exists", http.StatusUnauthorized)
			return
		}

		AccountStore[username] = password

		fmt.Fprintf(w, "Welcome %s, you can now login\n", username)
	})
}

func SignInHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if len(username) == 0 || len(password) == 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if password != AccountStore[username] {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		sessionId := uuid.New().String()

		CookieStore[sessionId] = username

		http.SetCookie(w, &http.Cookie{Name: "session-id", Value: sessionId})

		fmt.Fprintf(w, "Welcome %s, you are now logged in\n", username)
	})
}

func SignOutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, err := r.Cookie("session-id")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		username, ok := CookieStore[sessionId.Value]
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		delete(CookieStore, sessionId.Value)

		fmt.Fprintf(w, "See ya, %s\n", username)
	})
}

func MustAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, err := r.Cookie("session-id")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		username, ok := CookieStore[sessionId.Value]
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), nameKey, username)))
	})
}
