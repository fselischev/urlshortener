package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"urlshortener/models"
)

type Client struct {
	addr string
}

func New(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Shorten(r *models.ShortenReqest) (*models.Shorten, error) {
	data, _ := json.Marshal(r)

	resp, err := http.Post(c.addr+"/shorten", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var s *models.Shorten
	err = json.NewDecoder(resp.Body).Decode(&s)
	return s, err
}

func (c *Client) Go(r *models.GoReqest) error {
	resp, err := http.Get(c.addr + "/go/" + r.Key)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}
