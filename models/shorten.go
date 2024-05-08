package models

type ShortenReqest struct {
	URL string `json:"url"`
}

type GoReqest struct {
	Key string `json:"key"`
}

type Shorten struct {
	URL string `json:"url"`
	Key string `json:"key"`
}
