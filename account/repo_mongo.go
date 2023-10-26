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

type MongoRepository struct {
	client *mongo.Client
}

func (m *MongoRepository) Inject(client *mongo.Client) {
	m.client = client
}

func (m *MongoRepository) Create(ctx context.Context, account Account) error {
	collection := m.client.Database("instrumentality").Collection("accounts")
	_, err := collection.InsertOne(ctx, encode(account))
	return err
}

func (m *MongoRepository) FindByUsername(ctx context.Context, username Username) (Account, error) {
	collection := m.client.Database("instrumentality").Collection("accounts")
	var mongoAccount MongoAccount
	if err := collection.FindOne(ctx, bson.D{{"username", username}}).Decode(&mongoAccount); err != nil {
		return Account{}, errors.Wrap(err, "invalid username")
	}
	return decode(mongoAccount), nil
}

func (m *MongoRepository) FindByEmail(_ context.Context, _ Email) (Account, error) {
	panic("not implemented") // TODO: Implement
}

func (m *MongoRepository) createIndex() error {
	collection := m.client.Database("instrumentality").Collection("accounts")
	ctx := context.Background()
	model := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return err
	}

	return nil
}
