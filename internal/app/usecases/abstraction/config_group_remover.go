package abstraction

type ConfigGroupRemover interface {
	Remove(id int) error
}
