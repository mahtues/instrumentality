package account

import (
	"fmt"
	"net/http"

	"github.com/mahtues/form"
	"github.com/pkg/errors"
)

type Handler struct {
	mux     *http.ServeMux
	service *AccountImpl
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func NewHandler(service *AccountImpl) *Handler {
	return (&Handler{}).Init(service)
}

func (h *Handler) Init(service *AccountImpl) *Handler {
	h.mux = http.NewServeMux()

	h.mux.HandleFunc("/auth/signup", h.signUp)
	h.mux.HandleFunc("/auth/signin", h.signIn)

	h.service = service

	return h
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	MustMethodFunc(http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		frm, err := createFormFromRequest(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := h.service.Create(r.Context(), frm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Welcome, %s.\n", frm.Username)
	})(w, r)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	MustMethodFunc(http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		frm, err := verifyFormFromRequest(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := h.service.Verify(r.Context(), frm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Welcome back, %s.\n", frm.Username)
	})(w, r)
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

func MustMethodFunc(method string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		next(w, r)
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