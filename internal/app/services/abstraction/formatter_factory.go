package abstraction

type FormatterFactory interface {
	Fetch(formatterType string) (ConfigsFormatter, error)
}
