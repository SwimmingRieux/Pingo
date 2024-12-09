package repository

type ConfigRepository interface {
	CreateGroup(groupName string) (int, error)
	CreateConfig(groupId int, path string) (int, error)
}
