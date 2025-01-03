package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"pingo/configs"
	"pingo/internal/domain/repository"
)

type ConfigRemover struct {
	configRepository repository.ConfigRepository
	configuration    configs.Configuration
}

func (remover *ConfigRemover) Remove(id int) error {
	config, err := remover.configRepository.GetConfig(id)
	if err != nil {
		errText := remover.configuration.Errors.ConfigNotFound
		return fmt.Errorf("%v %w", errText, err)
	}
	defaultPath := remover.configuration.V2.ConfigurationPath
	filePath := filepath.Join(defaultPath, config.Path)

	if err = remover.configRepository.DeleteConfig(id); err != nil {
		errText := remover.configuration.Errors.FileRemoveError
		return fmt.Errorf("%v %v %w", errText, config.Path, err)
	}

	if err = os.Remove(filePath); err != nil {
		errText := remover.configuration.Errors.ConfigRemoveError
		return fmt.Errorf("%v %w", errText, err)
	}
	return nil
}
