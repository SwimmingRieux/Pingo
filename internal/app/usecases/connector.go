package usecases

import (
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/repository"
)

type Connector struct {
	configRepository repository.ConfigRepository
	configActivator  abstraction.ConfigActivator
	recorder         abstraction.NetworkLogRecorder
}

func (connector *Connector) Connect(configId int) error {
	config, err := connector.configRepository.GetConfig(configId)
	if err != nil {
		return err
	}
	if err = connector.configActivator.Activate(config.Path); err != nil {
		return err
	}
	go connector.recorder.Record()
	return nil
}

// todo: disconnector
