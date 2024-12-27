package repository

import (
	"pingo/internal/domain/dtos"
	"pingo/internal/domain/entities"
)

type ConfigRepository interface {
	CreateConfig(groupId int, path string, configType string) (int, error)
	DeleteConfig(id int) error
	GetConfig(id int) (entities.Config, error)
	UpdateConfig(id int, configDto dtos.UpdateConfigDto)
}
