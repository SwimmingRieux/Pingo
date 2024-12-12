package usecases

import (
	"fmt"
	"pingo/configs/abstraction"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/repository"
)

type ConfigGroupReceiver struct {
	repository   repository.GroupRepository
	configReader abstraction.Config
}

func (receiver *ConfigGroupReceiver) Get(id int) (entities.Group, error) {
	group, err := receiver.repository.GetGroup(id)
	if err != nil {
		errText, _ := receiver.configReader.Get("errors.group_not_found")
		return group, fmt.Errorf("%v %w", errText, err)
	}
	return group, nil
}
