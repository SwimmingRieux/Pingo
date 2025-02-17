package repository

type RepositoryConfigDeleter interface {
	DeleteConfig(id int) error
}
