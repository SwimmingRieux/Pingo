package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type JsonConfig struct {
	configMap map[string]interface{}
	mu        sync.RWMutex
	once      sync.Once
}

var (
	instance *JsonConfig
	once     sync.Once
)

func GetInstance() *JsonConfig {
	once.Do(func() {
		instance = &JsonConfig{}
	})
	return instance
}

func (c *JsonConfig) Initialize(path string) error {
	var err error
	c.once.Do(func() {
		jsonData, e := os.ReadFile(path)
		if e != nil {
			err = fmt.Errorf("error reading config file: %w", e)
			return
		}
		if e = json.Unmarshal(jsonData, &c.configMap); e != nil {
			err = fmt.Errorf("error unmarshaling config file: %w", e)
		}
	})
	return err
}

func (c *JsonConfig) Get(variablePath string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	parts := strings.Split(variablePath, ".")
	currentMap := c.configMap
	for _, part := range parts {
		var ok bool
		currentMap, ok = currentMap[part].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("config path %q not found", variablePath)
		}
	}
	return fmt.Sprintf("%v", currentMap), nil
}
