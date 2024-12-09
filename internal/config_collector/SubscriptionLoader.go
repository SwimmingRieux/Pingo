package config_collector

import (
	"errors"
	"io"
	"net/http"
)

type SubscriptionLoader struct {
	configs string
}

func (s SubscriptionLoader) GetSub(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Error: HTTP status " + string(rune(resp.StatusCode)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	text := string(body)
	s.configs = text
	return text, nil
}
