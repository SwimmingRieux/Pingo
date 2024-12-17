package usecases

import (
	"pingo/internal/domain/repository"
)

type ConfigsOrganizer struct {
	domainRepository repository.DomainRepository
	configRepository repository.ConfigRepository
}

func (organizer *ConfigsOrganizer) Organize(groupId int, domainsCountLimit int) error {
	// get some free ports from OS
	// keep the ports busy somehow, so OS don't give them to someone else
	// get all configs of groupId from configsRepo
	// set the port of each config from free ports
	// get top domainsCountLimit from domainsRepo
	// call goroutines of the pinger function which fill another channel
	// call a goroutine which reads pings of pairs from the second pair and stores in database
	// for configs for domains put every pair in first channel
	// when the whole bisiness is done it's over
}
