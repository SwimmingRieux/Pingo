package abstraction

type ConfigCreator interface {
	Create(input string) error
}
