package account

import (
	"fmt"
	"net/http"

	"github.com/mahtues/form"
	"github.com/pkg/errors"
)

func InitializeHandlers(mux *http.ServeMux, service Service) error {

	mux.Handle("/signup", signUpHandler(service))
	mux.Handle("/signin", signInHandler(service))

	return nil
}

func signUpHandler(service Service) http.Handler {
	return MustMethod(http.MethodPost, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		form, err := createFormFromRequest(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := service.Create(r.Context(), form); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Welcome, %s.\n", form.Username)
	}))
}

func signInHandler(service Service) http.Handler {
	return MustMethod(http.MethodPost, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		form, err := verifyFormFromRequest(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := service.Verify(r.Context(), form); err != nil {
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

func createFormFromRequest(r *http.Request) (CreateForm, error) {
	createForm := CreateForm{}

	if err := form.Unmarshal(r, &createForm); err != nil {
		return CreateForm{}, err
	}

	if len(createForm.Username) == 0 || len(createForm.Password) == 0 || len(createForm.Email) == 0 {
		return CreateForm{}, errors.New("missing username or password")
	}

	return createForm, nil
}

func verifyFormFromRequest(r *http.Request) (VerifyForm, error) {
	verifyForm := VerifyForm{}

	if err := form.Unmarshal(r, &verifyForm); err != nil {
		return VerifyForm{}, err
	}

	if len(verifyForm.Username) == 0 || len(verifyForm.Password) == 0 {
		return VerifyForm{}, errors.New("missing username or password")
	}

	return verifyForm, nil
}
