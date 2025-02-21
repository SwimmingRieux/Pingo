package services

import (
	"errors"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/structs"
	"strings"
)

type ConfigsCollectionFormatter struct {
	formatterFactory abstraction.FormatterFactory
	configuration    *configs.Configuration
}

func NewConfigsCollectionFormatter(formatterFactory abstraction.FormatterFactory, configuration *configs.Configuration) *ConfigsCollectionFormatter {
	return &ConfigsCollectionFormatter{
		formatterFactory: formatterFactory,
		configuration:    configuration,
	}
}

func (f *ConfigsCollectionFormatter) FormatCollection(rawConfigs []string) ([]structs.FormattedConfigAndType, error) {
	var formattedConfigs []structs.FormattedConfigAndType
	for _, rawConfig := range rawConfigs {
		configType := strings.Split(rawConfig, "://")[0]
		formatter, err := f.formatterFactory.Fetch(configType)
		if err != nil {
			continue
		}

		formattedConfig, err := formatter.Format(rawConfig)
		formattedConfigAndType := structs.FormattedConfigAndType{FormattedConfig: formattedConfig, Type: configType}
		if err == nil {
			formattedConfigs = append(formattedConfigs, formattedConfigAndType)
		}
	}
	if len(formattedConfigs) == 0 {
		errText := f.configuration.Errors.ConfigFormatError
		return nil, errors.New(errText)
	}
	return formattedConfigs, nil
}
