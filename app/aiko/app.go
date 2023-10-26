package main

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mahtues/instrumentality/account"
	"github.com/mahtues/instrumentality/zmisc"
)

type Config struct {
	Port        string
	MongoDbHost string
}

type App struct {
	config Config

	httpHandler *http.ServeMux
	mongoClient *mongo.Client

	accountService    account.Service
	accountRepository account.MongoRepository
	accountHandler    account.Handler
}

func NewApp(config Config) (*App, error) {
	app := &App{
		config:      config,
		httpHandler: http.NewServeMux(),
	}

	var err error

	// map handlers
	app.httpHandler.HandleFunc("/", app.helloHandlerFunc)

	// initialize resources
	app.mongoClient, err = zmisc.NewMongoClient(app.config.MongoDbHost)
	if err != nil {
		return nil, errors.WithMessage(err, "error creating client for mongodb")
	}

	// initialize components
	app.accountRepository.Inject(
		app.mongoClient,
	)

	app.accountService.Inject(
		&app.accountRepository,
	)

	app.accountHandler.Inject(
		"/auth",
		&app.accountService,
	)
	app.httpHandler.Handle("/auth/", &app.accountHandler)

	return app, nil
}

func (s *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.httpHandler.ServeHTTP(w, r)
}

func (s *App) helloHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from app running on port %v\n", s.config.Port)
}
