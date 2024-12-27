package abstraction

import (
	"net"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/structs"
	"sync"
)

type ConfigCollectionPinger interface {
	PingAllConfigs(configs []entities.Config, domainsWithRank []structs.DomainWithRank, wg *sync.WaitGroup, listeners []net.Listener, configScoresMap *sync.Map)
}
