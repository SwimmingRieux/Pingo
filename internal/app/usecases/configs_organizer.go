package usecases

import (
	"fmt"
	"net"
	configurations "pingo/configs"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/dtos"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/repository"
	"pingo/internal/domain/structs"
	"sync"
)

type ConfigsOrganizer struct {
	domainRepository  repository.DomainRepository
	configRepository  repository.ConfigRepository
	groupRepository   repository.GroupRepository
	portSetterFactory abstraction.PortSetterFactory
	pinger            abstraction.ConfigPinger
	configuration     configurations.Configuration
}

func (organizer *ConfigsOrganizer) Organize(groupId int, domainsCountLimit int) error {

	configs, err := organizer.getConfigs(groupId)
	if err != nil {
		return err
	}

	listeners, err := organizer.getListeners(configs)
	if err != nil {
		return err
	}
	defer organizer.closeAllListeners(listeners)

	if err = organizer.setPortOnConfigs(configs, listeners); err != nil {
		return err
	}

	domainsWithRank, err := organizer.getDomainsWithRank(domainsCountLimit)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var configScoresMap sync.Map
	organizer.pingAllConfigs(configs, domainsWithRank, &wg, listeners, &configScoresMap)
	wg.Wait()

	organizer.writeScoresToDb(configs, &configScoresMap)

	return nil
}

func (organizer *ConfigsOrganizer) writeScoresToDb(configs []entities.Config, configScoresMap *sync.Map) {
	for _, config := range configs {
		value, ok := configScoresMap.Load(config)
		if ok {
			configDto := dtos.UpdateConfigDto{Type: config.Type, Path: config.Path, Score: value.(float64)}
			organizer.configRepository.UpdateConfig(config.ConfigId, configDto)
		} else {
			configDto := dtos.UpdateConfigDto{Type: config.Type, Path: config.Path, Score: 0}
			organizer.configRepository.UpdateConfig(config.ConfigId, configDto)
		}
	}
}

func (organizer *ConfigsOrganizer) pingAllConfigs(configs []entities.Config, domainsWithRank []structs.DomainWithRank, wg *sync.WaitGroup, listeners []net.Listener, configScoresMap *sync.Map) {
	maxGoroutines := organizer.configuration.PingerGoroutinesMax
	semaphore := make(chan struct{}, maxGoroutines)

	for i, config := range configs {
		for _, domain := range domainsWithRank {
			wg.Add(1)
			go func() {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				organizer.pinger.Ping(config, domain, listeners[i], configScoresMap)
			}()
		}
	}
}

func (organizer *ConfigsOrganizer) getDomainsWithRank(domainsCountLimit int) ([]structs.DomainWithRank, error) {
	domains, err := organizer.domainRepository.GetDomains(domainsCountLimit)
	if err != nil {
		return nil, err
	}

	domainsCount := len(domains)
	domainsWithRank := make([]structs.DomainWithRank, domainsCount)

	for i, domain := range domains {
		domainsWithRank[i] = structs.DomainWithRank{Domain: domain, Rank: i}
	}
	return domainsWithRank, nil
}

func (organizer *ConfigsOrganizer) setPortOnConfigs(configs []entities.Config, listeners []net.Listener) error {
	for i, config := range configs {
		portSetter, err := organizer.portSetterFactory.Fetch(config.Type)
		if err != nil {
			return err
		}
		portSetter.SetPort(listeners[i], config.Path)
	}
	return nil
}

func (organizer *ConfigsOrganizer) getListeners(configs []entities.Config) ([]net.Listener, error) {
	listenerIterationLimit := organizer.configuration.ListenerIterationLimit
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
		errText := configurations.Errors.NotEnoughPortsFound
		return nil, fmt.Errorf("%v", errText)
	}
	return listeners, nil
}

func (organizer *ConfigsOrganizer) closeAllListeners(listeners []net.Listener) {
	for _, listener := range listeners {
		listener.Close()
	}
}

func (organizer *ConfigsOrganizer) getConfigs(groupId int) ([]entities.Config, error) {
	configs, err := organizer.groupRepository.GetConfigs(groupId)
	if err != nil {
		errText := organizer.configuration.Errors.GroupNotFound
		return nil, fmt.Errorf("%v %v", errText, groupId)
	}
	return configs, nil
}
