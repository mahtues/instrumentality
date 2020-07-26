package account

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type mongoModel struct {
	Username string `bson:"username"`
	Hash     string `bson:"hash"`
	Email    string `bson:"email"`
}

func getClient(ctx context.Context) (*mongo.Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongodb:27017"))
}

func getCollection(ctx context.Context) (*mongo.Collection, error) {
	client, err := getClient(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database("instrumentality").Collection("accounts"), err
}

type mongoService struct {
}

func (s *mongoService) Create(ctx context.Context, form CreateForm) error {
	accounts, err := getCollection(ctx)
	if err != nil {
		log.Println("connect error:", err.Error())
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	account := mongoModel{Username: form.Username, Hash: string(hash), Email: form.Email}

	_, err = accounts.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true)})
	if err != nil {
		return err
	}

	_, err = accounts.InsertOne(ctx, account)

	return err
}

func (s *mongoService) Verify(ctx context.Context, form VerifyForm) error {
	accounts, err := getCollection(ctx)
	if err != nil {
		log.Println("connect error:", err.Error())
		return err
	}

	filter := bson.D{{"username", form.Username}}

	var account mongoModel
	if err := accounts.FindOne(ctx, filter).Decode(&account); err != nil {
		return errors.Wrap(err, "invalid username")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Hash), []byte(form.Password))
	if err != nil {
		return errors.Wrap(err, "invalid password")
	}

	return nil
}
