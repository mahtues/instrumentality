package concrete

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/mahtues/instrumentality/internal/account"
)

type AccountImpl struct {
	repository Repository
}

func New(mongo *mongo.Client) (*AccountImpl, error) {
	repository, err := NewMongoRepository(mongo)
	if err != nil {
		return nil, errors.Wrap(err, "error creating repository")
	}

	return &AccountImpl{
		repository: repository,
	}, nil
}

func (s *AccountImpl) Create(ctx context.Context, form account.CreateForm) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	acc := account.Account{
		Username: account.Username(form.Username),
		Hash:     account.Hash(hash),
		Email:    account.Email(form.Email),
	}

	return s.repository.Create(ctx, acc)
}

func (s *AccountImpl) Verify(ctx context.Context, form account.VerifyForm) error {
	var acc account.Account
	var err error

	if acc, err = s.repository.FindByUsername(ctx, account.Username(form.Username)); err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Hash), []byte(form.Password))
	if err != nil {
		return errors.Wrap(err, "invalid password")
	}

	return nil
}
