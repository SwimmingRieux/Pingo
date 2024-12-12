package abstraction

type LogReceiver interface {
	GetDomains(limit int) ([]string, error)
}
