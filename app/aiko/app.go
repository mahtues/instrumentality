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

	handler *http.ServeMux

	mongo *mongo.Client

	accountService    account.Service
	accountRepository account.MongoRepository
	accountHandler    account.Handler
}

func NewApp(config Config) (*App, error) {
	app := &App{
		config:  config,
		handler: http.NewServeMux(),
	}

	// initialize resources
	if err := app.initResources(); err != nil {
		return nil, errors.Wrap(err, "error initializing resources")
	}

	// initialize adapters
	app.accountRepository.Inject(
		app.mongo,
	)

	// initialize services
	app.accountService.Inject(
		&app.accountRepository,
	)

	// initialize handlers
	app.accountHandler.Inject(
		"/auth",
		&app.accountService,
	)

	// map handlers
	app.handler.Handle("/auth/", &app.accountHandler)

	app.handler.HandleFunc("/", app.helloHandlerFunc)

	return app, nil
}

func (s *App) initResources() error {
	var err error

	s.mongo, err = zmisc.NewMongoClient(s.config.MongoDbHost)
	if err != nil {
		return errors.Wrap(err, "error creating client for mongodb")
	}

	return nil
}

func (s *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *App) helloHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from app running on port %v\n", s.config.Port)
}
