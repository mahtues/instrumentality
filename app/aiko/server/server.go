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

type Server struct {
	config Config

	handler *http.ServeMux

	mongo *mongo.Client

	accountService    account.Service
	accountRepository account.MongoRepository
	accountHandler    account.Handler
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

	// initialize adapters
	server.accountRepository.Inject(server.mongo)

	// initialize services
	server.accountService.Inject(&server.accountRepository)

	// initialize handlers
	server.accountHandler.Inject("/auth", &server.accountService)

	// map sub routes
	server.handler.HandleFunc("/", server.helloHandlerFunc)
	server.handler.Handle("/auth/", &server.accountHandler)

	return server, nil
}

func (s *Server) initResources() error {
	var err error

	s.mongo, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_HOST")).SetMonitor(apmmongo.CommandMonitor()))
	if err != nil {
		return errors.Wrap(err, "error creating client for mongodb")
	}

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) helloHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from app running on port %d\n", s.config.Port)
}
