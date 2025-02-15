package repository

import "pingo/internal/domain/entities"

type RepositoryConfigRetriever interface {
	GetConfig(id int) (entities.Config, error)
}
