package repository

import "pingo/internal/domain/entities"

type ConfigRepository interface {
	CreateConfig(groupId int, path string, configType string) (int, error)
	DeleteConfig(id int) error
	GetConfig(id int) (entities.Config, error)
}
