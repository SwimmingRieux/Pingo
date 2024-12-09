package abstract

type ISubscriptionLoader interface {
	GetSub(url string) (string, error)
}
