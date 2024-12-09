package entities

type Config struct {
	ConfigId int `gorm:"primarykey"`
	Path     string
	Score    float64
	GroupID  int
	Group    Group `gorm:"foreignKey:GroupId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
}
