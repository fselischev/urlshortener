//go:build !solution

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"urlshortener/pkg/shortener"
)

func main() {
	shortener := shortener.New()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", shortener.ShortenHandler)
	mux.HandleFunc("GET /go/", shortener.GoHandler)

	portFlag := flag.Int("port", 6029, "port number")
	flag.Parse()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), mux); err != nil {
		log.Printf("Failed to start server: %v\n", err)
	}
}
