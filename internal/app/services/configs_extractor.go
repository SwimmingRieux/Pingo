package services

import (
	"encoding/base64"
	"strings"
	"time"
)

type ConfigsExtractor struct{}

func (extractor *ConfigsExtractor) Extract(configs string) (string, []string) {
	if checkEncoding(configs) {
		decoded, err := base64.StdEncoding.DecodeString(configs)
		if err != nil {
			return "", nil
		}
		configs = string(decoded)
	}

	splitConfigs := strings.Fields(configs)

	validConfigs := []string{}
	for _, part := range splitConfigs {
		if strings.HasPrefix(part, "vless://") || strings.HasPrefix(part, "vmess://") || strings.HasPrefix(part, "trojan://") || strings.HasPrefix(part, "ss://") {
			validConfigs = append(validConfigs, part)
		}
	}

	groupName := time.Now().Format("Jan 02 2006 15:04")

	return groupName, validConfigs
}

func checkEncoding(input string) bool {
	if len(input) != 4 {
		return false
	}
	for _, ch := range input {
		if !(ch == '+' || (ch >= '/' && ch <= '9') || ch == '=' || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')) {
			return false
		}
	}
	return true
}
