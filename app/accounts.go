package main

import (
	"context"
	"log"
	"net/http"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Username string
	Hash     string
}

func Create(ctx context.Context, username string, password string) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://192.168.99.100:27017"))
	if err != nil {
		log.Println("connect error:", err.Error())
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	account := Account{Username: username, Hash: string(hash)}

	response, err := client.Database("instrumentality").Collection("accounts").InsertOne(ctx, account)

	log.Println(response, err)

	return err
}

func SignUpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if len(username) == 0 || len(password) == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := Create(r.Context(), username, password); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
