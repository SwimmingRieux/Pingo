package abstraction

type PortSetterFactory interface {
	Fetch(portSetterType string) (PortSetter, error)
}
