package account

import (
	"context"
	"errors"
	"log"
	"net/http"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"
)

type AccountModel struct {
	Username string
	Hash     string
}

type AccountForm struct {
	Username string
	Password string
}

func formFromRequest(r *http.Request) (AccountForm, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) == 0 || len(password) == 0 {
		return AccountForm{}, errors.New("missing username or password")
	}

	return AccountForm{Username: username, Password: password}, nil
}

func Create(ctx context.Context, form AccountForm) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://192.168.99.100:27017"))
	if err != nil {
		log.Println("connect error:", err.Error())
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	account := AccountModel{Username: form.Username, Hash: string(hash)}

	response, err := client.Database("instrumentality").Collection("accounts").InsertOne(ctx, account)

	log.Println(response, err)

	return err
}

func SignUpHandler() http.Handler {
	return MustMethod(http.MethodPost, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		form, err := formFromRequest(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := Create(r.Context(), form); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
