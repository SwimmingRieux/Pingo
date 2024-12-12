package usecases

import (
	"pingo/internal/domain/repository"
)

type LogReceiver struct {
	domainRepository repository.DomainRepository
}

func (receiver *LogReceiver) GetDomains(limit int) ([]string, error) {
	domains, err := receiver.domainRepository.GetDomains(limit)
	if err != nil {
		return []string{}, err
	}
	domainStrings := make([]string, len(domains))
	for i, d := range domains {
		domainStrings[i] = d.Address
	}
	return domainStrings, nil
}
