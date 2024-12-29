package abstraction

type Connector interface {
	Connect(configId int) error
}
