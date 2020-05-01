package api

import (
	"fmt"
	"net/http"

	"github.com/kingnido/instrumentality/account"
	"github.com/pkg/errors"
)

func createFormFromRequest(r *http.Request) (account.CreateForm, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if len(username) == 0 || len(password) == 0 || len(email) == 0 {
		return account.CreateForm{}, errors.New("missing username or password")
	}

	return account.CreateForm{Username: username, Password: password, Email: email}, nil
}

func verifyFormFromRequest(r *http.Request) (account.VerifyForm, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) == 0 || len(password) == 0 {
		return account.VerifyForm{}, errors.New("missing username or password")
	}

	return account.VerifyForm{Username: username, Password: password}, nil
}

func SignUpHandler() http.Handler {
	return MustMethod(http.MethodPost, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		form, err := createFormFromRequest(r)
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
		form, err := verifyFormFromRequest(r)
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
