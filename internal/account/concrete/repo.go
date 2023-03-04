package concrete

import (
	"context"

	"github.com/mahtues/instrumentality/internal/account"
)

type Repository interface {
	Create(context.Context, account.Account) error
	FindByUsername(context.Context, account.Username) (account.Account, error)
	FindByEmail(context.Context, account.Email) (account.Account, error)
}
