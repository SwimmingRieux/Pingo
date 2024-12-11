package repository

type GroupRepository interface {
	CreateGroup(groupName string) (int, error)
}
