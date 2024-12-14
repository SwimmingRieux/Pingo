package services

import (
	"fmt"
	configAbstraction "pingo/configs/abstraction"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/enums"
	"sync"
)

type FormatterFactory struct {
	formatters   sync.Map
	configReader configAbstraction.Config
}

func (factory *FormatterFactory) Fetch(formatterType enums.FormatterType) (abstraction.ConfigsFormatter, error) {
	if cached, ok := factory.formatters.Load(formatterType); ok {
		return cached.(abstraction.ConfigsFormatter), nil
	}
	var formatter abstraction.ConfigsFormatter

	switch formatterType {
	case enums.VmessFormatter:
		formatter = &VmessConfigsFormatter{}
	case enums.VlessFormatter:
		formatter = &VlessConfigsFormatter{}
	case enums.TrojanFormatter:
		formatter = &TrojanConfigsFormatter{}
	case enums.SsFormatter:
		formatter = &SsConfigsFormatter{}
	default:
		errText, _ := factory.configReader.Get("invalid_formatter")
		return nil, fmt.Errorf("%v %v", errText, formatterType)
	}

	factory.formatters.Store(formatterType, formatter)
	return formatter, nil
}
