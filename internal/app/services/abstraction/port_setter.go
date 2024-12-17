package abstraction

type PortSetter interface {
	SetPort(port int, configPath string) error
}
