package usecases

import (
	"errors"
	"fmt"
	"os"
	"path"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/repository"
	"pingo/internal/domain/structs"
	"strings"
	"sync"
)

type ConfigCreator struct {
	loader           abstraction.UrlLoader
	extractor        abstraction.ConfigsExtractor
	collectionWriter abstraction.ConfigCollectionFileWriter
	groupRepository  repository.RepositoryGroupCreator
	configuration    configs.Configuration
	formatterFactory abstraction.FormatterFactory
}

func (creator *ConfigCreator) Create(input string) error {
	var configsString string
	var err error

	if strings.HasPrefix(input, "http") && !strings.Contains(input, " ") {
		configsString, err = creator.loader.Load(input)
		if err != nil {
			errText := creator.configuration.Errors.LoadFromLinkError
			return fmt.Errorf("%v %w", errText, err)
		}
	} else {
		configsString = input
	}

	groupName, rawConfigs := creator.extractor.Extract(configsString)
	if len(rawConfigs) == 0 {
		errText := creator.configuration.Errors.ConfigNotFound
		return errors.New(errText)
	}

	var formattedConfigs []structs.FormattedConfigAndType
	for _, rawConfig := range rawConfigs {
		configType := strings.Split(rawConfig, "://")[0]
		formatter, err := creator.formatterFactory.Fetch(configType)
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
		errText := creator.configuration.Errors.ConfigFormatError
		return errors.New(errText)
	}

	newGroupId, err := creator.groupRepository.CreateGroup(groupName)
	if err != nil {
		errText := creator.configuration.Errors.GroupCreatingError
		return fmt.Errorf("%v %w", errText, err)
	}

	pingoPath := os.Getenv("PINGO_PATH")
	v2ConfigsPath := creator.configuration.V2.ConfigurationPath

	groupPath := path.Join(pingoPath, v2ConfigsPath, groupName)
	err = os.Mkdir(groupPath, 0755)
	if err != nil {
		errText := creator.configuration.Errors.DirectoryCreatingError
		return fmt.Errorf("%v %w", errText, err)
	}

	var wg sync.WaitGroup
	creator.collectionWriter.WriteConfigsToFiles(formattedConfigs, &wg, groupPath, newGroupId)
	wg.Wait()

	return nil
}
