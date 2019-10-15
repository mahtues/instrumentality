package main

import (
	"errors"

	"github.com/google/uuid"
)

type InMemorySessionStore struct {
	store map[string]string
}

func NewInMemorySessionStore() InMemorySessionStore {
	return InMemorySessionStore{
		store: map[string]string{},
	}
}

func (ss *InMemorySessionStore) Add(username string) (string, error) {
	sessionId := uuid.New().String()
	ss.store[sessionId] = username
	return sessionId, nil
}

func (ss *InMemorySessionStore) Get(sessionId string) (string, error) {
	username, ok := ss.store[sessionId]
	if !ok {
		return "", errors.New("session not found")
	}

	return username, nil
}

func (ss *InMemorySessionStore) Del(sessionId string) (username string, err error) {
	username, ok := ss.store[sessionId]
	if !ok {
		return "", errors.New("session not found")
	}

	delete(ss.store, sessionId)

	return username, nil
}

var SessionStore = NewInMemorySessionStore()
