package services

import (
	"fmt"
	"os"
	"pingo/configs"
)

type ConfigsWriter struct {
	configuration configs.Configuration
}

func (writer *ConfigsWriter) Write(jsonConfig string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		errText := writer.configuration.Errors.FileCreatingError
		return fmt.Errorf("%v %w", errText, err)
	}
	defer file.Close()

	_, err = file.WriteString(jsonConfig)
	if err != nil {
		errText := writer.configuration.Errors.WriteToFileError
		return fmt.Errorf("%v %w", errText, err)
	}

	return nil
}
