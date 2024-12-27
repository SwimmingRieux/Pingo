package abstraction

import (
	"net"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/structs"
	"sync"
)

type ConfigPinger interface {
	Ping(config entities.Config, domain structs.DomainWithRank,
		listener net.Listener, domainScoresMap *sync.Map)
}
