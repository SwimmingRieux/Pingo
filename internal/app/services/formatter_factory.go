package services

import (
	"fmt"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"sync"
)

type FormatterFactory struct {
	formatters    sync.Map
	configuration *configs.Configuration
}

func NewFormatterFactory(configuration *configs.Configuration) *FormatterFactory {
	return &FormatterFactory{
		configuration: configuration,
		formatters:    sync.Map{},
	}
}

func (factory *FormatterFactory) Fetch(formatterType string) (abstraction.ConfigsFormatter, error) {
	if cached, ok := factory.formatters.Load(formatterType); ok {
		return cached.(abstraction.ConfigsFormatter), nil
	}
	var formatter abstraction.ConfigsFormatter

	switch formatterType {
	case "vmess":
		formatter = NewVmessConfigsFormatter(factory.configuration)
	case "vless":
		formatter = NewVlessConfigsFormatter(factory.configuration)
	case "trojan":
		formatter = NewTrojanConfigsFormatter(factory.configuration)
	case "ss":
		formatter = NewSsConfigsFormatter(factory.configuration)
	default:
		errText := factory.configuration.Errors.InvalidFormatter
		return nil, fmt.Errorf("%v %v", errText, formatterType)
	}

	factory.formatters.Store(formatterType, formatter)
	return formatter, nil
}
