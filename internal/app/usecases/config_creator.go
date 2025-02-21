package usecases

import (
	"errors"
	"fmt"
	"os"
	"path"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/repository"
	"strings"
	"sync"
)

type ConfigCreator struct {
	loader              abstraction.UrlLoader
	extractor           abstraction.ConfigsExtractor
	collectionWriter    abstraction.ConfigCollectionFileWriter
	groupRepository     repository.RepositoryGroupCreator
	configuration       configs.Configuration
	collectionFormatter abstraction.ConfigsCollectionFormatter
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

	formattedConfigs, err := creator.collectionFormatter.FormatCollection(rawConfigs)
	if err != nil {
		errText := creator.configuration.Errors.CollectiveFormatError
		return fmt.Errorf("%v %w", errText, err)
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
