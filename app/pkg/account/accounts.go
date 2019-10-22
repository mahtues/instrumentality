package account

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors"
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

func getClient(ctx context.Context) (*mongo.Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI("mongodb://192.168.99.100:27017"))
}

func getCollection(ctx context.Context) (*mongo.Collection, error) {
	client, err := getClient(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database("instrumentality").Collection("accounts"), err
}

func Create(ctx context.Context, form AccountForm) error {
	accounts, err := getCollection(ctx)
	if err != nil {
		log.Println("connect error:", err.Error())
		return err
	}
	defer accounts.Client().Disconnect()

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	account := AccountModel{Username: form.Username, Hash: string(hash)}

	_, err = accounts.InsertOne(ctx, account)

	return err
}

func Verify(ctx context.Context, form AccountForm) error {
	accounts, err := getCollection(ctx)
	if err != nil {
		log.Println("connect error:", err.Error())
		return err
	}

	var account AccountModel
	if err := accounts.FindOne(ctx, bson.M{"username": form.Username}).Decode(&account); err != nil {
		return errors.Wrap(err, "invalid username")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Hash), []byte(form.Password))
	if err != nil {
		return errors.Wrap(err, "invalid password")
	}

	return nil
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

		if err := Verify(r.Context(), form); err != nil {
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
