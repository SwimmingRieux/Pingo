package usecases

import (
	"context"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/repository"
)

type Connector struct {
	configRepository  repository.RepositoryConfigRetriever
	configActivator   abstraction.ConfigActivator
	configDeactivator abstraction.ConfigDeactivator
	recorder          abstraction.NetworkLogRecorder
	cancelFunc        context.CancelFunc
	killFunc          func() error
}

func (connector *Connector) Connect(configId int) error {
	connector.Disconnect()

	config, err := connector.configRepository.GetConfig(configId)
	if err != nil {
		return err
	}
	if connector.killFunc, err = connector.configActivator.Activate(config.Path); err != nil {
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
	connector.configDeactivator.Deactivate(connector.killFunc)
}
