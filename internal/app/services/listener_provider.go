package services

import (
	"fmt"
	"net"
	"pingo/configs"
)

type ListenerProvider struct {
	configuration *configs.Configuration
}

func NewListenerProvider(configuration *configs.Configuration) *ListenerProvider {
	return &ListenerProvider{configuration: configuration}
}

func (s *ListenerProvider) GetListeners(configsLength int) ([]net.Listener, error) {
	listenerIterationLimit := s.configuration.ListenerIterationLimit
	var listeners []net.Listener
	counter := 0
	for len(listeners) < configsLength && counter < listenerIterationLimit {
		listener, err := net.Listen("tcp", ":0")
		if err == nil {
			listeners = append(listeners, listener)
		}
		counter++
	}
	if len(listeners) != configsLength {
		errText := s.configuration.Errors.NotEnoughPortsFound
		return nil, fmt.Errorf("%v", errText)
	}
	return listeners, nil
}
