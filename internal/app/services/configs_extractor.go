package services

import (
	"encoding/base64"
	"strings"
	"time"
)

type ConfigsExtractor struct{}

func (extractor *ConfigsExtractor) Extract(configs string) (string, []string) {
	decoded, err := base64.StdEncoding.DecodeString(configs)

	if err == nil {
		configs = string(decoded)
	}

	splitConfigs := strings.Fields(configs)

	validConfigs := []string{}
	for _, part := range splitConfigs {
		if strings.HasPrefix(part, "vless") || strings.HasPrefix(part, "vmess") || strings.HasPrefix(part, "trojan") || strings.HasPrefix(part, "shadowsocks") {
			validConfigs = append(validConfigs, part)
		}
	}

	groupName := time.Now().Format("Jan 02 2006 15:04")

	return groupName, validConfigs
}
