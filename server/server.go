package server

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mahtues/instrumentality/account"
)

type Config struct {
	Port int
}

type Services struct {
	Account *account.AccountImpl
}

type Resources struct {
	Mongo *mongo.Client
}

type Server struct {
	Config Config

	Handler *http.ServeMux

	Resources Resources

	Services Services
}

func NewServer(config Config) (*Server, error) {
	server := &Server{
		Config:  config,
		Handler: http.NewServeMux(),
	}

	// initialize resources
	if err := server.InitResources(); err != nil {
		return nil, errors.Wrap(err, "error initializing resources")
	}

	// initialize services
	if err := server.InitServices(); err != nil {
		return nil, errors.Wrap(err, "error initializing services")
	}

	if err := server.InjectServices(); err != nil {
		return nil, errors.Wrap(err, "error injecting services into other services")
	}

	// initialize handlers
	if err := server.InitHandlers(); err != nil {
		return nil, errors.Wrap(err, "error initializing handlers")
	}

	return server, nil
}

func (s *Server) InitResources() error {
	var err error

	s.Resources.Mongo, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_HOST")).SetMonitor(apmmongo.CommandMonitor()))
	if err != nil {
		return errors.Wrap(err, "error creating client for mongodb")
	}

	return nil
}

func (s *Server) InitServices() error {
	var err error

	s.Services.Account, err = account.New(s.Resources.Mongo)
	if err != nil {
		return errors.Wrap(err, "error initializing account service")
	}

	return nil
}

func (s Server) InjectServices() error {
	return nil
}

func (s Server) InitHandlers() error {
	s.Handler.HandleFunc("/", s.helloHandlerFunc)

	s.Handler.Handle("/auth/", account.NewHandler(s.Services.Account))

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}

func (s *Server) helloHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from app running on port %d\n", s.Config.Port)
}
