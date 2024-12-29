package abstraction

type ConfigActivator interface {
	Activate(path string) error
}
