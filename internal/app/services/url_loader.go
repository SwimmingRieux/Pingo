package services

import (
	"fmt"
	"io"
	"net/http"
	"pingo/configs/abstraction"
)

type UrlLoader struct {
	configReader abstraction.Config
}

func (loader *UrlLoader) Load(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errText, _ := loader.configReader.Get("errors.http_status")
		return "", fmt.Errorf("%v %v", errText, resp.StatusCode)

	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	text := string(body)
	return text, nil

}
