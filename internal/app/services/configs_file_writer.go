package services

import (
	"fmt"
	"os"
	configAbstraction "pingo/configs/abstraction"
)

type ConfigsWriter struct {
	configReader configAbstraction.Config
}

func (writer *ConfigsWriter) Write(jsonConfig string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		errText, _ := writer.configReader.Get("file_creating_error")
		return fmt.Errorf("%v", errText)
	}
	defer file.Close()

	_, err = file.WriteString(jsonConfig)
	if err != nil {
		errText, _ := writer.configReader.Get("write_to_file_error")
		return fmt.Errorf("%v %v", errText, err)
	}

	return nil
}
