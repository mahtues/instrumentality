package api

import (
	"fmt"
	"net/http"

	"github.com/kingnido/instrumentality/app/pkg/account"
	"github.com/pkg/errors"
)

func formFromRequest(r *http.Request) (account.AccountForm, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) == 0 || len(password) == 0 {
		return account.AccountForm{}, errors.New("missing username or password")
	}

	return account.AccountForm{Username: username, Password: password}, nil
}

func SignUpHandler() http.Handler {
	return MustMethod(http.MethodPost, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		form, err := formFromRequest(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := account.Create(r.Context(), form); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Welcome, %s.\n", form.Username)
	}))
}

func SignInHandler() http.Handler {
	return MustMethod(http.MethodPost, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		form, err := formFromRequest(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := account.Verify(r.Context(), form); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Welcome back, %s.\n", form.Username)
	}))
}

func MustMethod(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}
