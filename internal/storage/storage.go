package storage

import (
	"context"
	"database/sql"
	"os"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
}

func New(ctx context.Context) (*Storage, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("PGCONN"))
	if err != nil {
		return nil, err
	}

	if _, err := conn.Exec(ctx, `
	CREATE TABLE url_mappings (
            key VARCHAR(7) PRIMARY KEY,
            url TEXT NOT NULL UNIQUE
        );
        CREATE INDEX idx_key ON url_mappings (key);
        CREATE INDEX idx_url ON url_mappings (url);	
	`); err != nil {
		return nil, err
	}

	return &Storage{db: conn}, nil
}

func (s *Storage) Insert(ctx context.Context, url, key string) error {
	if _, err := s.db.Exec(ctx,
		`INSERT INTO url_mappings (key, url) VALUES ($1, $2)`,
		key, url); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetURL(ctx context.Context, key string) (string, error) {
	var url string
	if err := s.db.QueryRow(ctx,
		`SELECT url FROM url_mappings WHERE key = $1`,
		key).Scan(&url); err != nil {
		return "", sql.ErrNoRows
	}
	return url, nil
}

func (s *Storage) GetKey(ctx context.Context, url string) (string, error) {
	var key string
	if err := s.db.QueryRow(ctx,
		`SELECT key FROM url_mappings WHERE url = $1`,
		url).Scan(&key); err != nil {
		return "", sql.ErrNoRows
	}
	return key, nil
}

func (s *Storage) Close() error {
	return s.db.Close(context.Background())
}
