package abstraction

import "net"

type PortSetter interface {
	SetPort(listener net.Listener, configPath string) error
}
