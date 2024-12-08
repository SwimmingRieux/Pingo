package entities

type Domain struct {
	DomainId int `gorm:"primarykey"`
	Address  string
}
