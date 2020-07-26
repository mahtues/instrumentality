package account

import (
	"context"
)

type CreateForm struct {
	Username string
	Password string
	Email    string
}

type Creater interface {
	Create(context.Context, CreateForm) error
}

type VerifyForm struct {
	Username string
	Password string
}

type Verifier interface {
	Verify(context.Context, VerifyForm) error
}

// usefull to mocks in tests
type CreaterFunc func(context.Context, CreateForm) error

func (f CreaterFunc) Create(ctx context.Context, form CreateForm) error {
	return f(ctx, form)
}

type VerifierFunc func(context.Context, VerifyForm) error

func (f VerifierFunc) Verify(ctx context.Context, form VerifyForm) error {
	return f(ctx, form)
}

// a bunch of stuff to use the default service (temporary)
var defaultService = &mongoService{}

func DefaultCreater() Creater {
	return defaultService
}

func DefaultVerifier() Verifier {
	return defaultService
}

func Create(ctx context.Context, form CreateForm) error {
	return defaultService.Create(ctx, form)
}

func Verify(ctx context.Context, form VerifyForm) error {
	return defaultService.Verify(ctx, form)
}
