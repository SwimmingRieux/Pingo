package abstraction

type ConfigsWriter interface {
	Write(jsonConfig string, path string) error
}
