package abstraction

type ConfigsExtractor interface {
	Extract(input string) (string, []string)
}
