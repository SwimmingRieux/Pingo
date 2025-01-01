package abstraction

import (
	"net"
	"pingo/internal/domain/entities"
)

type ListenerProvider interface {
	GetListeners(configs []entities.Config) ([]net.Listener, error)
}
