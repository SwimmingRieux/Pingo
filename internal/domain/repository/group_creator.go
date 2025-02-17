package repository

type RepositoryGroupCreator interface {
	CreateGroup(groupName string) (int, error)
}
