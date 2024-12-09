package entities

type Group struct {
	GroupId int `gorm:"primarykey"`
	Name    string
	Configs []Config
}
