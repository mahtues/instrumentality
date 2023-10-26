package server

import (
	"fmt"
	"log"
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
	log.Print("initializing resources")
	if err := server.initResources(); err != nil {
		return nil, errors.Wrap(err, "error initializing resources")
	}
	log.Print("resources initialized")

	// initialize adapters
	server.accountRepository.Inject(
		server.mongo,
	)

	// initialize services
	server.accountService.Inject(
		&server.accountRepository,
	)

	// initialize handlers
	server.accountHandler.Inject(
		"/auth",
		&server.accountService,
	)

	// map handlers
	server.handler.Handle("/auth/", &server.accountHandler)

	server.handler.HandleFunc("/", server.helloHandlerFunc)

	return server, nil
}

func (s *Server) initResources() error {
	var err error

	s.mongo, err = zmisc.NewMongoClient(s.config.MongoDbHost)
	if err != nil {
		return errors.Wrap(err, "error creating client for mongodb")
	}

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) helloHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from app running on port %v\n", s.config.Port)
}
