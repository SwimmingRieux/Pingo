package services

import (
	"fmt"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"sync"
)

type FormatterFactory struct {
	formatters    sync.Map
	configuration configs.Configuration
}

func (factory *FormatterFactory) Fetch(formatterType string) (abstraction.ConfigsFormatter, error) {
	if cached, ok := factory.formatters.Load(formatterType); ok {
		return cached.(abstraction.ConfigsFormatter), nil
	}
	var formatter abstraction.ConfigsFormatter

	switch formatterType {
	case "vmess":
		formatter = &VmessConfigsFormatter{}
	case "vless":
		formatter = &VlessConfigsFormatter{}
	case "trojan":
		formatter = &TrojanConfigsFormatter{}
	case "ss":
		formatter = &SsConfigsFormatter{}
	default:
		errText := factory.configuration.Errors.InvalidFormatter
		return nil, fmt.Errorf("%v %v", errText, formatterType)
	}

	factory.formatters.Store(formatterType, formatter)
	return formatter, nil
}
