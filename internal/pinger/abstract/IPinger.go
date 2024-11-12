package abstract

type IPinger interface {
	GetAveragePing(port int) (int, error)
}
