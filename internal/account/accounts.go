package account

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

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

type VerifyForm struct {
	Username string
	Password string
}

var defaultRepository, repoErr = NewMongoRepository(context.Background())

func Create(ctx context.Context, form CreateForm) error {
	if defaultRepository == nil {
		return errors.Wrap(repoErr, "missing repository")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	account := Account{Username: Username(form.Username), Hash: Hash(hash), Email: Email(form.Email)}

	return defaultRepository.Create(ctx, account)
}

func Verify(ctx context.Context, form VerifyForm) error {
	if defaultRepository == nil {
		return errors.Wrap(repoErr, "missing repository")
	}

	var account Account
	var err error

	if account, err = defaultRepository.FindByUsername(ctx, Username(form.Username)); err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Hash), []byte(form.Password))
	if err != nil {
		return errors.Wrap(err, "invalid password")
	}

	return nil
}

type Service interface {
}
