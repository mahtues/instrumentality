package main

// imports
import (
	"fmt"
	"log"

	"go.elastic.co/apm/module/apmhttp"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	port := 80

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), apmhttp.Wrap(mux)); err != nil {
		log.Fatal(err)
	}
}
