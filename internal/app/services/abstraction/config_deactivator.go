package abstraction

type ConfigDeactivator interface {
	Deactivate(killFunc func() error) error
}
