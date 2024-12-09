package abstraction

type UrlLoader interface {
	Load(url string) (string, error)
}
