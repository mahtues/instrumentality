package main

// testezera

import (
	"net/http"
)

type StringError string

func (se StringError) Error() string {
	return string(se)
}

const SessionAlreadyExistsErr = StringError("session already exists")

type Session struct {
	Id     string
	Values map[string]string
}

const SessionKey = "instrumentality-session-id"

type Store interface {
	// create a new session, or return a existing one and fail
	NewSession(r *http.Request) (*Session, error)

	// returns current session, or fail if not exists
	GetSession(r *http.Request) (*Session, error)

	// save session
	SaveSession(r *http.Request, w http.ResponseWriter, s *Session)
}

type InMemoryStore struct {
	Sessions map[string]*Session
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Sessions: map[string]*Session{},
	}
}

func (ims *InMemoryStore) NewSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(SessionKey)
	if err == nil {
		sessionId := cookie.Value
		session, ok := ims.Sessions[sessionId]
		if ok {
			return session, SessionAlreadyExistsErr
		}
	}

	return nil, nil
}
