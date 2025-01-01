package repository

import "pingo/internal/domain/entities"

type DomainRepository interface {
	GetDomains(limit int) ([]entities.Domain, error)
	AddDomains(domains []string)
}
