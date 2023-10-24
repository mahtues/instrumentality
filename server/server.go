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

// todo: make this one private?
type Config struct {
	Port int
}

type services struct {
	account *account.Service
}

type resources struct {
	mongo *mongo.Client
}

type Server struct {
	config    Config
	handler   *http.ServeMux
	resources resources
	services  services
}

func NewServer(config Config) (*Server, error) {
	server := &Server{
		config:  config,
		handler: http.NewServeMux(),
	}

	// initialize resources
	if err := server.initResources(); err != nil {
		return nil, errors.Wrap(err, "error initializing resources")
	}

	// initialize services
	if err := server.initServices(); err != nil {
		return nil, errors.Wrap(err, "error initializing services")
	}

	if err := server.injectServices(); err != nil {
		return nil, errors.Wrap(err, "error injecting services into other services")
	}

	// initialize handlers
	if err := server.initHandlers(); err != nil {
		return nil, errors.Wrap(err, "error initializing handlers")
	}

	return server, nil
}

func (s *Server) initResources() error {
	var err error

	s.resources.mongo, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_HOST")).SetMonitor(apmmongo.CommandMonitor()))
	if err != nil {
		return errors.Wrap(err, "error creating client for mongodb")
	}

	return nil
}

func (s *Server) initServices() error {
	var err error

	s.services.account, err = account.New(s.resources.mongo)
	if err != nil {
		return errors.Wrap(err, "error initializing account service")
	}

	return nil
}

func (s *Server) injectServices() error {
	return nil
}

func (s *Server) initHandlers() error {
	s.handler.HandleFunc("/", s.helloHandlerFunc)

	s.handler.Handle("/auth/", account.NewHandler(s.services.account))

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) helloHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from app running on port %d\n", s.config.Port)
}
