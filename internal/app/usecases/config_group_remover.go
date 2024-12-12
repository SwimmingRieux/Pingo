package usecases

import (
	"pingo/internal/domain/repository"
)

type ConfigGroupRemover struct {
	groupRepository repository.GroupRepository
}

func (remover *ConfigGroupRemover) Remove(id int) error {
	if err := remover.groupRepository.DeleteGroup(id); err != nil {
		return err
	}
	return nil
}
