package abstraction

import "pingo/internal/domain/structs"

type DomainRankFetcher interface {
	GetDomainsWithRank(domainsCountLimit int) ([]structs.DomainWithRank, error)
}
