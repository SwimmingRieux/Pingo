package abstraction

type ConfigRemover interface {
	Remove(id int) error
}
