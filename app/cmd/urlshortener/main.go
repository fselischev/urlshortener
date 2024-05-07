package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"urlshortener/internal/storage"
	"urlshortener/pkg/shortener"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	endpointRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "url_shortener_endpoint_requests_total",
		Help: "Total number of requests to each endpoint",
	}, []string{"endpoint", "client"})

	endpointDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "url_shortener_endpoint_duration_seconds",
		Help:    "Duration of requests to each endpoint",
		Buckets: prometheus.DefBuckets,
	}, []string{"endpoint"})
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
	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime).Seconds()
			endpointDuration.WithLabelValues("/shorten").Observe(duration)
		}()

		endpointRequests.WithLabelValues("/shorten", "client").Inc()
		shortener.ShortenHandler(w, r)
	})

	mux.HandleFunc("/go/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime).Seconds()
			endpointDuration.WithLabelValues("/go/").Observe(duration)
		}()

		endpointRequests.WithLabelValues("/go/", "client").Inc()
		shortener.GoHandler(w, r)
	})

	mux.Handle("/metrics", promhttp.Handler())

	portFlag := flag.Int("port", 6029, "port number")
	flag.Parse()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), mux); err != nil {
		log.Printf("Failed to start server: %v\n", err)
	}
}
