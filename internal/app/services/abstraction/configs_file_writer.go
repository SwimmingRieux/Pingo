package abstraction

type ConfigFileWriter interface {
	Write(jsonConfig string, path string) error
}
