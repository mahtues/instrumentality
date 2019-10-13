package service

import (
	"fmt"
)

type Authorization interface {
	SignUp(SignUpForm) error
	SignIn(SignInForm) error
	SignOut(SignOutForm) error
}

type SignUpForm struct{}

type SignInForm struct{}

type SignOutForm struct{}

func None() {
	x := make([]int)
	x = append(x, 5)
	fmt.Println(len(x))
}
