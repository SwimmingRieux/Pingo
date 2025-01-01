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
	configuration configs.Configuration
}

func (s *ConfigCollectionPinger) PingAllConfigs(configs []entities.Config, domainsWithRank []structs.DomainWithRank, wg *sync.WaitGroup, listeners []net.Listener, configScoresMap *sync.Map) {
	maxGoroutines := s.configuration.PingerGoroutinesMax
	semaphore := make(chan struct{}, maxGoroutines)

	for i, config := range configs {
		for _, domain := range domainsWithRank {
			wg.Add(1)
			go func() {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				s.pinger.Ping(config, domain, listeners[i], configScoresMap)
			}()
		}
	}
}
