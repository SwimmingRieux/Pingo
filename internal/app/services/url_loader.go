package services

import (
	"fmt"
	"io"
	"net/http"
	"pingo/configs"
)

type UrlLoader struct {
	configuration configs.Configuration
}

func (loader *UrlLoader) Load(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errText := loader.configuration.Errors.HttpStatus
		return "", fmt.Errorf("%v %v", errText, resp.StatusCode)

	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	text := string(body)
	return text, nil

}
