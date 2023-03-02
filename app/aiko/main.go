package main

import (
	"fmt"
	"log"
	"os"

	"context"
	"net/http"

	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmhttp"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mahtues/instrumentality/internal/account"
)

func main() {
	config := Config{
		Port: 8080,
	}

	server, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), apmhttp.Wrap(server)); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Port int
}

type Services struct {
	Account account.Service
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

func NewServer(config Config) (Server, error) {
	server := Server{
		Config:  config,
		Handler: http.NewServeMux(),
	}

	// initialize resources
	if err := server.InitResources(); err != nil {
		return Server{}, errors.Wrap(err, "error initializing resources")
	}

	// intialize services

	server.Handler.Handle("/", server.helloHandler())
	server.Handler.Handle("/signup", account.SignUpHandler())
	server.Handler.Handle("/signin", account.SignInHandler())

	return server, nil
}

func (s Server) InitResources() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_HOST")).SetMonitor(apmmongo.CommandMonitor()))
	if err != nil {
		return errors.Wrap(err, "error creating client for mongodb")
	}
	s.Resources.Mongo = client

	return nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}

func (s Server) helloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from app running on port %d\n", s.Config.Port)
	})
}
