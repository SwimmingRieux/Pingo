package services

import (
	"pingo/internal/domain/repository"
	"pingo/internal/domain/structs"
)

type DomainRankFetcher struct {
	domainRepository repository.DomainRepository
}

func (s *DomainRankFetcher) GetDomainsWithRank(domainsCountLimit int) ([]structs.DomainWithRank, error) {
	domains, err := s.domainRepository.GetDomains(domainsCountLimit)
	if err != nil {
		return nil, err
	}

	domainsCount := len(domains)
	domainsWithRank := make([]structs.DomainWithRank, domainsCount)

	for i, domain := range domains {
		domainsWithRank[i] = structs.DomainWithRank{Domain: domain, Rank: i}
	}
	return domainsWithRank, nil
}
