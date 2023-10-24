package account

import (
	"context"
)

type IRepository interface {
	Create(context.Context, Account) error
	FindByUsername(context.Context, Username) (Account, error)
	FindByEmail(context.Context, Email) (Account, error)
}
