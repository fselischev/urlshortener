package main

import (
	"log"
	"urlshortener/client"
	"urlshortener/models"
)

func main() {
	cl := client.New("http://localhost:6029")

	sh, err := cl.Shorten(&models.ShortenReqest{
		URL: "http://ifmo.ru",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("url: %v, key: %v", sh.URL, sh.Key)

	if err := cl.Go(&models.GoReqest{
		Key: sh.Key,
	}); err != nil {
		log.Fatal(err)
	}
}
