package shortener

import (
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"urlshortener/internal/storage"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type URLShortener struct {
	l *zap.SugaredLogger

	db *storage.Storage
}

func New(db *storage.Storage) *URLShortener {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapConfig.EncoderConfig.TimeKey = ""
	zapConfig.EncoderConfig.EncodeTime = nil
	zapConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	zapConfig.EncoderConfig.EncodeCaller = nil
	zapConfig.Encoding = "console"
	zapConfig.OutputPaths = []string{"stdout"}
	l, _ := zapConfig.Build()
	defer func() { _ = l.Sync() }()

	return &URLShortener{
		l:  l.Sugar().Named("UrlShortener"),
		db: db,
	}
}

func (us *URLShortener) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		us.l.Errorf("invalid reqest: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	defer func() {
		r.Body.Close()
		_, _ = io.Copy(io.Discard, r.Body)
	}()

	ctx := context.Background()
	url := requestData.URL

	key, err := us.db.GetKey(ctx, url)
	if key != "" && err != nil {
		us.l.Errorf("cannot get key in db: %v", err)

		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	if key == "" {
		key = generateRandomKey()
		if err := us.db.Insert(ctx, url, key); err != nil {
			us.l.Errorf("cannot insert in db: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(
		struct {
			URL string `json:"url"`
			Key string `json:"key"`
		}{
			URL: requestData.URL,
			Key: key,
		},
	)

	us.l.Infof("url: %v, key: %v", url, key)
}

func (us *URLShortener) GoHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/go/"):]
	ctx := context.Background()
	url, err := us.db.GetURL(ctx, key)
	if err != nil {
		us.l.Errorf("cannot get url from db: %v", err)
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
	us.l.Infof("redirected %v using key %v", url, key)
}

func generateRandomKey() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 7)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
