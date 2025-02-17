package repository

import "pingo/internal/domain/entities"

type RepositoryGroupRetriever interface {
	GetGroup(id int) (entities.Group, error)
}
