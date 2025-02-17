package repository

import "pingo/internal/domain/entities"

type RepositoryDomainRetriever interface {
	GetDomains(limit int) ([]entities.Domain, error)
}
