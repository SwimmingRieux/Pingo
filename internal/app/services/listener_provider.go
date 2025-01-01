package services

import (
	"fmt"
	"net"
	"pingo/configs"
	"pingo/internal/domain/entities"
)

type ListenerProvider struct {
	configuration configs.Configuration
}

func (s *ListenerProvider) GetListeners(configs []entities.Config) ([]net.Listener, error) {
	listenerIterationLimit := s.configuration.ListenerIterationLimit
	var listeners []net.Listener
	counter := 0
	for len(listeners) < len(configs) && counter < listenerIterationLimit {
		listener, err := net.Listen("tcp", ":0")
		if err == nil {
			listeners = append(listeners, listener)
		}
		counter++
	}
	if len(listeners) != len(configs) {
		errText := s.configuration.Errors.NotEnoughPortsFound
		return nil, fmt.Errorf("%v", errText)
	}
	return listeners, nil
}
