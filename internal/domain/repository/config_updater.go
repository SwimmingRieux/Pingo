package repository

import "pingo/internal/domain/dtos"

type RepositoryConfigUpdater interface {
	UpdateConfig(id int, configDto dtos.UpdateConfigDto)
}
