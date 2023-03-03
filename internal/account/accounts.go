package account

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository Repository
}

func New(mongo *mongo.Client) (Service, error) {
	repository, err := NewMongoRepository(mongo)
	if err != nil {
		return Service{}, errors.Wrap(err, "error creating repository")
	}

	return Service{
		repository: repository,
	}, nil
}

type Username string
type Hash string
type Email string

type Account struct {
	Username Username
	Hash     Hash
	Email    Email
}

type CreateForm struct {
	Username string `form:"Username"`
	Password string `form:"Password"`
	Email    string `form:"Email"`
}

func (s Service) Create(ctx context.Context, form CreateForm) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	account := Account{Username: Username(form.Username), Hash: Hash(hash), Email: Email(form.Email)}

	return s.repository.Create(ctx, account)
}

type VerifyForm struct {
	Username string
	Password string
}

func (s Service) Verify(ctx context.Context, form VerifyForm) error {
	var account Account
	var err error

	if account, err = s.repository.FindByUsername(ctx, Username(form.Username)); err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Hash), []byte(form.Password))
	if err != nil {
		return errors.Wrap(err, "invalid password")
	}

	return nil
}
