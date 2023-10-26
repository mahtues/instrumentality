package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"go.elastic.co/apm/module/apmhttp"

	"github.com/mahtues/instrumentality/zmisc"
)

var (
	appPortArg     = flag.String("app-port", "", "app listening port")
	appPortEnv     = os.Getenv("APP_PORT")
	appPortDefault = "8080"

	mongoDbHostArg     = flag.String("mongodb-host", "", "mongodb host")
	mongoDbHostEnv     = os.Getenv("MONGODB_HOST")
	mongoDbHostDefault = "mongodb://localhost:27017"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	flag.Parse()

	var (
		appPort     = zmisc.FirstNonEmpty(*appPortArg, appPortEnv, appPortDefault)
		mongoDbHost = zmisc.FirstNonEmpty(*mongoDbHostArg, mongoDbHostEnv, mongoDbHostDefault)
	)

	config := Config{
		Port:        appPort,
		MongoDbHost: mongoDbHost,
	}

	log.Printf("app config: %#v", config)

	app, err := NewApp(config)
	if err != nil {
		log.Fatal(err)
	}

	srv := http.Server{
		Addr:     fmt.Sprintf(":%v", config.Port),
		Handler:  apmhttp.Wrap(app),
		ErrorLog: log.Default(),
	}

	// setting up signal capturing
	stop := make(chan os.Signal, 1)
	closed := make(chan struct{}, 1)

	signal.Notify(stop, os.Interrupt)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				log.Printf("waiting for server to shutdown")
				if err := srv.Shutdown(context.TODO()); err != nil {
					log.Print(err)
				}
			case <-closed:
				return
			}
		}
	}()

	log.Printf("starting server on port %v", config.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	closed <- struct{}{}

	wg.Wait()

	log.Printf("server closed")
}
