package services

import (
	"encoding/base64"
	"strings"
	"time"
)

type ConfigsExtractor struct{}

func (extractor *ConfigsExtractor) Extract(configs string) (string, []string) {
	checkEncoding := CheckEncoding(configs)
	if checkEncoding {
		decoded, err := base64.StdEncoding.DecodeString(configs)
		if err != nil {
			return "", nil
		} else {
			configs = string(decoded)
		}
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

func CheckEncoding(input string) bool {
	for _, ch := range input {
		if !(ch == 43 || (ch >= 47 && ch <= 57) || ch == 61 || (ch >= 65 && ch <= 90) || (ch >= 97 && ch <= 122)) {
			return false // config wasn't already base64 encoded
		}
	}
	return true // config was already base64 encoded
}
