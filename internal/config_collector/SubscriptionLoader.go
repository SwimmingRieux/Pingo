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

	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Error: HTTP status " + string(rune(resp.StatusCode)))
	}

	// Read Context Of Subscription Page
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//Convert Body To String
	text := string(body)
	s.configs = text
	return text, nil
}
