package account

import (
	"context"
	"log"

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
