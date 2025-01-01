package usecases

import (
	"context"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/repository"
)

type Connector struct {
	configRepository repository.ConfigRepository
	configActivator  abstraction.ConfigActivator
	recorder         abstraction.NetworkLogRecorder
	cancelFunc       context.CancelFunc
}

func (connector *Connector) Connect(configId int) error {
	connector.Disconnect()

	config, err := connector.configRepository.GetConfig(configId)
	if err != nil {
		return err
	}
	if err = connector.configActivator.Activate(config.Path); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	connector.cancelFunc = cancel
	go connector.recorder.Record(ctx)
	return nil
}

func (connector *Connector) Disconnect() {
	if connector.cancelFunc != nil {
		connector.cancelFunc()
	}
}
