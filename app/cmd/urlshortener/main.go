//go:build !solution

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"urlshortener/internal/storage"
	"urlshortener/pkg/shortener"
)

func main() {
	ctx := context.Background()
	db, err := storage.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	shortener := shortener.New(db)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", shortener.ShortenHandler)
	mux.HandleFunc("GET /go/", shortener.GoHandler)

	portFlag := flag.Int("port", 6029, "port number")
	flag.Parse()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), mux); err != nil {
		log.Printf("Failed to start server: %v\n", err)
	}
}
