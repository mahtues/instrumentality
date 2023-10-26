package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.elastic.co/apm/module/apmhttp"

	"github.com/mahtues/instrumentality/zmisc"
)

func main() {
	var (
		appPortArg     = flag.String("app-port", "", "app listening port")
		appPortEnv     = os.Getenv("APP_PORT")
		appPortDefault = "8080"

		mongoDbHostArg     = flag.String("mongodb-host", "", "mongodb host")
		mongoDbHostEnv     = os.Getenv("MONGODB_HOST")
		mongoDbHostDefault = "mongodb://mongodb:27017"
	)
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	var (
		appPort     = zmisc.FirstNonEmpty(*appPortArg, appPortEnv, appPortDefault)
		mongoDbHost = zmisc.FirstNonEmpty(*mongoDbHostArg, mongoDbHostEnv, mongoDbHostDefault)
	)

	config := Config{
		Port:        appPort,
		MongoDbHost: mongoDbHost,
	}

	log.Printf("app will listen on port %v", config.Port)

	app, err := NewApp(config)
	if err != nil {
		log.Fatal(err)
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%v", config.Port),
		Handler: apmhttp.Wrap(app),
	}

	log.Printf("app listening on port %v", config.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
