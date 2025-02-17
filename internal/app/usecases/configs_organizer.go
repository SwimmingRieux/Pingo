package usecases

import (
	"fmt"
	"net"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/repository"
	"sync"
)

type ConfigsOrganizer struct {
	groupRepository repository.RepositoryConfigsRetriever

	portSetterFactory abstraction.PortSetterFactory
	configuration     configs.Configuration

	configScoreWriter abstraction.ConfigScoreWriter
	configPinger      abstraction.ConfigCollectionPinger
	domainRankFetcher abstraction.DomainRankFetcher
	listenerProvider  abstraction.ListenerProvider
}

func (organizer *ConfigsOrganizer) Organize(groupId int, domainsCountLimit int) error {

	configs, err := organizer.getConfigs(groupId)
	if err != nil {
		return err
	}

	listeners, err := organizer.listenerProvider.GetListeners(len(configs))
	if err != nil {
		return err
	}
	defer organizer.closeAllListeners(listeners)

	if err = organizer.setPortOnConfigs(configs, listeners); err != nil {
		return err
	}

	domainsWithRank, err := organizer.domainRankFetcher.GetDomainsWithRank(domainsCountLimit)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var configScoresMap sync.Map
	organizer.configPinger.PingAllConfigs(configs, domainsWithRank, &wg, listeners, &configScoresMap)
	wg.Wait()

	organizer.configScoreWriter.WriteScoresToDb(configs, &configScoresMap)

	return nil
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
