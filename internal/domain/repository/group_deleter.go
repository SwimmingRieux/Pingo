package repository

type RepositoryGroupDeleter interface {
	DeleteGroup(id int) error
}
