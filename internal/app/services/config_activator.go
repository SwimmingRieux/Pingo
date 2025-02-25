package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Config struct {
	Inbounds []Inbound `json:"inbounds"`
}

type Inbound struct {
	Listen string `json:"listen"`
	Port   int    `json:"port"`
	Tag    string `json:"tag"`
}

type ConfigActivator struct {
}

func (c *ConfigActivator) Activate(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return fmt.Errorf("could not parse config file: %w", err)
	}

	var ports []int
	for _, inbound := range config.Inbounds {
		ports = append(ports, inbound.Port)
	}

	for _, port := range ports {
		os.Setenv("HTTP_PROXY", fmt.Sprintf("127.0.0.1:%d", port))
		os.Setenv("HTTPS_PROXY", fmt.Sprintf("127.0.0.1:%d", port))
		os.Setenv("FTP_PROXY", fmt.Sprintf("127.0.0.1:%d", port))
		os.Setenv("SOCKS_HOST", fmt.Sprintf("127.0.0.1:%d", port))
	}

	cmd := exec.Command("v2ray", "-c", path)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run V2Ray: %w", err)
	}

	return nil
}
