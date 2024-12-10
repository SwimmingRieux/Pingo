package repository

type ConfigRepository interface {
	CreateConfig(groupId int, path string) (int, error)
}
