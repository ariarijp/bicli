package main

import (
	"fmt"
	"strings"
)

type BicliConfig struct {
	Login  string `toml:"login"`
	APIKey string `toml:"api-key"`
}

type BitlyResponse struct {
	StatusCode int               `json:"status_code"`
	StatusText string            `json:"status_text"`
	Data       BitlyResponseData `json:"data"`
}

type BitlyResponseData struct {
	LongURL    string `json:"long_url"`
	URL        string `json:"url"`
	Hash       string `json:"hash"`
	GlobalHash string `json:"global_hash"`
	NewHash    int    `json:"new_hash"`
}

type BitlyErrorResponse struct {
	StatusCode int                    `json:"status_code"`
	StatusText string                 `json:"status_text"`
	Data       BitlyResponseErrorData `json:"data"`
}

type BitlyResponseErrorData struct {
}

type ShortURL struct {
	LineNum int
	URL     string
	LongURL string
}

func (s ShortURL) ToCSV(sep string) string {
	return strings.Join([]string{
		fmt.Sprintf("%d", s.LineNum),
		s.URL,
		s.LongURL,
	}, sep)
}

type ShortURLs []ShortURL
