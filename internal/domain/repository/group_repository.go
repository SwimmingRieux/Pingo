package repository

import "pingo/internal/domain/entities"

type GroupRepository interface {
	CreateGroup(groupName string) (int, error)
	GetGroup(id int) (entities.Group, error)
	DeleteGroup(id int) error
}
