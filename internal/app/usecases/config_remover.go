package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"pingo/configs/abstraction"
	"pingo/internal/domain/repository"
)

type ConfigRemover struct {
	configRepository repository.ConfigRepository
	configReader     abstraction.Config
}

func (remover *ConfigRemover) Remove(id int) error {
	config, err := remover.configRepository.GetConfig(id)
	if err != nil {
		errText, _ := remover.configReader.Get("errors.config_not_found")
		return fmt.Errorf("%v %w", errText, err)
	}
	defaultPath, _ := remover.configReader.Get("v2.config_path")
	filePath := filepath.Join(defaultPath, config.Path)

	if err = remover.configRepository.DeleteConfig(id); err != nil {
		errText, _ := remover.configReader.Get("file_remove_error")
		return fmt.Errorf("%v %v %w", errText, config.Path, err)
	}

	if err = os.Remove(filePath); err != nil {
		errText, _ := remover.configReader.Get("config_remove_error")
		return fmt.Errorf("%v %w", errText, err)
	}
	return nil
}
