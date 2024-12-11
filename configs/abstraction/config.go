package abstraction

type Config interface {
	Initialize(path string) error
	Get(variablePath string) (string, error)
}
