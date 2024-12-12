package usecases

import (
	"pingo/internal/domain/entities"
	"pingo/internal/domain/repository"
)

type ConfigGroupReceiver struct {
	repository repository.GroupRepository
}

func (receiver *ConfigGroupReceiver) Get(id int) (entities.Group, error) {
	group, err := receiver.repository.GetGroup(id)
	return group, err
}
