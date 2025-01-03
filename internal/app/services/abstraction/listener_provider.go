package abstraction

import (
	"net"
)

type ListenerProvider interface {
	GetListeners(configsLength int) ([]net.Listener, error)
}
