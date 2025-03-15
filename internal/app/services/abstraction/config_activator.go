package abstraction

type ConfigActivator interface {
	Activate(path string) (func() error, error)
}
