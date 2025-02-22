package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"pingo/configs"
	"pingo/internal/domain/repository"
)

type ConfigRemover struct {
	repositoryConfigDeleter   repository.RepositoryConfigDeleter
	repositoryConfigRetriever repository.RepositoryConfigRetriever
	configuration             *configs.Configuration
}

func NewConfigRemover(deleter repository.RepositoryConfigDeleter, retriever repository.RepositoryConfigRetriever, configuration *configs.Configuration) *ConfigRemover {
	return &ConfigRemover{
		repositoryConfigDeleter:   deleter,
		repositoryConfigRetriever: retriever,
		configuration:             configuration,
	}
}

func (remover *ConfigRemover) Remove(id int) error {
	config, err := remover.repositoryConfigRetriever.GetConfig(id)
	if err != nil {
		errText := remover.configuration.Errors.ConfigNotFound
		return fmt.Errorf("%v %w", errText, err)
	}
	pingoPath := os.Getenv("PINGO_PATH")
	filePath := filepath.Join(pingoPath, remover.configuration.V2.ConfigurationPath, config.Path)

	if err = remover.repositoryConfigDeleter.DeleteConfig(id); err != nil {
		errText := remover.configuration.Errors.ConfigRemoveError
		return fmt.Errorf("%v %v %w", errText, config.Path, err)
	}

	if err = os.Remove(filePath); err != nil {
		errText := remover.configuration.Errors.FileRemoveError
		return fmt.Errorf("%v %w", errText, err)
	}
	return nil
}
