package handler

import (
	"fmt"
	"net/http"

	"go.elastic.co/apm/module/apmmongo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDb() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client, err := mongo.Connect(r.Context(), options.Client().ApplyURI("mongodb://mongodb:27017").SetMonitor(apmmongo.CommandMonitor()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := client.Ping(r.Context(), nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "oook")
	})
}
