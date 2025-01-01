package structs

import "pingo/internal/domain/entities"

type DomainWithRank struct {
	Domain entities.Domain
	Rank   int
}
