package account

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository Repository
}

func (s *Service) Inject(repository Repository) {
	s.repository = repository
}

func (s *Service) Create(ctx context.Context, form CreateForm) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	acc := Account{
		Username: Username(form.Username),
		Hash:     Hash(hash),
		Email:    Email(form.Email),
	}

	return s.repository.Create(ctx, acc)
}

func (s *Service) Verify(ctx context.Context, form VerifyForm) error {
	var acc Account
	var err error

	if acc, err = s.repository.FindByUsername(ctx, Username(form.Username)); err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Hash), []byte(form.Password))
	if err != nil {
		return errors.Wrap(err, "invalid password")
	}

	return nil
}
