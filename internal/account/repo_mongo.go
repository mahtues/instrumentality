package account

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoAccount struct {
	Username string `bson:"username"`
	Hash     string `bson:"hash"`
	Email    string `bson:"email"`
}

type MongoRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client) (*MongoRepository, error) {
	database := client.Database("instrumentality")
	collection := database.Collection("accounts")

	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true)})
	if err != nil {
		return nil, err
	}

	return &MongoRepository{client, database, collection}, nil
}

func (m *MongoRepository) Create(ctx context.Context, account Account) error {
	_, err := m.collection.InsertOne(ctx, encode(account))
	return err
}

func (m *MongoRepository) FindByUsername(ctx context.Context, username Username) (Account, error) {
	var mongoAccount MongoAccount
	if err := m.collection.FindOne(ctx, bson.D{{"username", username}}).Decode(&mongoAccount); err != nil {
		return Account{}, errors.Wrap(err, "invalid username")
	}
	return decode(mongoAccount), nil
}

func (m *MongoRepository) FindByEmail(_ context.Context, _ Email) (Account, error) {
	panic("not implemented") // TODO: Implement
}

func encode(account Account) MongoAccount {
	return MongoAccount{
		Username: string(account.Username),
		Hash:     string(account.Hash),
		Email:    string(account.Email),
	}
}

func decode(mongoAccount MongoAccount) Account {
	return Account{
		Username: Username(mongoAccount.Username),
		Hash:     Hash(mongoAccount.Hash),
		Email:    Email(mongoAccount.Email),
	}
}
