package services

import (
	"net"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/structs"
	"sync"

	"pingo/internal/domain/entities"
)

type ConfigCollectionPinger struct {
	pinger        abstraction.ConfigPinger
	configuration *configs.Configuration
}

func NewConfigCollectionPinger(pinger abstraction.ConfigPinger, configuration *configs.Configuration) *ConfigCollectionPinger {
	configCollectionPinger := &ConfigCollectionPinger{
		pinger:        pinger,
		configuration: configuration,
	}
	return configCollectionPinger
}

func (s *ConfigCollectionPinger) PingAllConfigs(configs []entities.Config, domainsWithRank []structs.DomainWithRank, wg *sync.WaitGroup, listeners []net.Listener, configScoresMap *sync.Map) {

	maxGoroutines := s.configuration.PingerGoroutinesMax
	semaphore := make(chan struct{}, maxGoroutines)

	for i, config := range configs {
		for _, domain := range domainsWithRank {
			wg.Add(1)
			go func(config entities.Config, domain structs.DomainWithRank, listener net.Listener) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				s.pinger.Ping(config, domain, listener, configScoresMap)
			}(config, domain, listeners[i])
		}
	}
}
