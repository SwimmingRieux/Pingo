package repository

import "pingo/internal/domain/entities"

type RepositoryConfigsRetriever interface {
	GetConfigs(id int) ([]entities.Config, error)
}
