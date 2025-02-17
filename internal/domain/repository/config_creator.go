package repository

type RepositoryConfigCreator interface {
	CreateConfig(groupId int, path string, configType string) (int, error)
}
