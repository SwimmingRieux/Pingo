package services

import (
	"fmt"
	"io"
	"net/http"
)

type UrlLoader struct{}

func (loader *UrlLoader) Load(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: HTTP status %v", resp.StatusCode)

	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	text := string(body)
	return text, nil

}
