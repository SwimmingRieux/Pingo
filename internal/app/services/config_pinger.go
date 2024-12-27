package services

import (
	"net"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/structs"
	"sync"
)

type ConfigPinger struct {
}

func (pinger *ConfigPinger) Ping(config entities.Config, domain structs.DomainWithRank,
	listener net.Listener, domainScoresMap *sync.Map) {
	// don't forget to kill the v2ray process so it doesn't remain running
	// also don't forget to stop listening, before running the v2ray process
}
