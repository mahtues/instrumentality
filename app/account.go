package main

import (
	"errors"
)

type InMemoryAccountStore struct {
	store map[string]string
}

func NewInMemoryAccountStore() InMemoryAccountStore {
	return InMemoryAccountStore{
		store: map[string]string{},
	}
}

func (ac *InMemoryAccountStore) AddUser(username string, password string) error {
	if len(username) == 0 && len(password) == 0 {
		return errors.New("empty parameters")
	}

	ac.store[username] = password

	return nil
}

func (ac *InMemoryAccountStore) HasUser(username string) bool {
	_, ok := ac.store[username]
	return ok
}

func (ac *InMemoryAccountStore) CheckUser(username string, password string) error {
	if len(username) == 0 && len(password) == 0 {
		return errors.New("empty parameters")
	}

	if password != ac.store[username] {
		return errors.New("invalid username or password")
	}

	return nil
}

func (ac *InMemoryAccountStore) DelUser(username string, password string) error {
	if err := ac.CheckUser(username, password); err != nil {
		return err
	}

	delete(ac.store, username)

	return nil
}

var AccountStore = NewInMemoryAccountStore()
