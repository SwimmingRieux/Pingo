package services

import (
	"fmt"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"sync"
)

type PortSetterFactory struct {
	portSetters   sync.Map
	configuration *configs.Configuration
}

func NewPortSetterFactory(configuration *configs.Configuration) *PortSetterFactory {
	return &PortSetterFactory{
		portSetters:   sync.Map{},
		configuration: configuration,
	}
}

func (factory *PortSetterFactory) Fetch(portSetterType string) (abstraction.PortSetter, error) {
	if cached, ok := factory.portSetters.Load(portSetterType); ok {
		return cached.(abstraction.PortSetter), nil
	}

	var portSetter abstraction.PortSetter

	switch portSetterType {
	case "vmess":
		portSetter = &VmessTypeSetter{}
	case "vless":
		portSetter = &VlessTypeSetter{}
	case "trojan":
		portSetter = &TrojanTypeSetter{}
	case "ss":
		portSetter = &SsTypeSetter{}
	default:
		errText := factory.configuration.Errors.InvalidPortSetter
		return nil, fmt.Errorf("%v %v", errText, portSetterType)
	}

	factory.portSetters.Store(portSetterType, portSetter)
	return portSetter, nil
}
