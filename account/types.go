package account

import (
	"context"
)

type IService interface {
	Create(ctx context.Context, form CreateForm) error
	Verify(ctx context.Context, form VerifyForm) error
}

type CreateForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Email    string `form:"email"`
}

type VerifyForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type Username string
type Hash string
type Email string

type Account struct {
	Username Username
	Hash     Hash
	Email    Email
}
