package services

import (
	"encoding/base64"
	"strings"
	"time"
)

type ConfigsExtractor struct{}

func (extractor *ConfigsExtractor) Extract(configs string) (string, []string) {
	configs, _ = decodeBase64IfNeeded(configs)

	splitConfigs := strings.Fields(configs)

	validConfigs := make([]string, 0)
	for _, part := range splitConfigs {
		if strings.HasPrefix(part, "vless://") || strings.HasPrefix(part, "vmess://") || strings.HasPrefix(part, "trojan://") || strings.HasPrefix(part, "ss://") {
			validConfigs = append(validConfigs, part)
		}
	}

	groupName := time.Now().Format("Jan 02 2006 15:04")

	return groupName, validConfigs
}

func decodeBase64IfNeeded(b64string string) (string, error) {
	b64string = strings.TrimSpace(b64string)

	padding := len(b64string) % 4
	b64stringFix := b64string
	if padding != 0 {
		b64stringFix += "===="[:4-padding]
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(b64stringFix)

	if err != nil {
		return b64string, err
	}

	return string(decodedBytes), nil
}
