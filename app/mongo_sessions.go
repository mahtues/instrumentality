package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Ping() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://192.168.99.100:27017"))
	if err != nil {
		log.Println("connect error:", err.Error())
		return err
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("ping error:", err.Error())
		return err
	}

	database := client.Database("tmp")
	collection := database.Collection("coltmp")
	response, err := collection.InsertOne(ctx, bson.M{})

	log.Println(response, err)

	return nil
}

func PingHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := Ping(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "ping to mondodb went fiiiiiine\n")
	})
}
