run:
	docker compose up -d --wait && env PGCONN="host=127.0.0.1 port=5432 database=pg user=me password=pass" go run app/cmd/urlshortener/main.go 

down:
	docker compose down
