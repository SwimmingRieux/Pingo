package usecases

import (
	"errors"
	"fmt"
	"os"
	"path"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/repository"
	"strconv"
	"strings"
	"sync"
)

type ConfigCreator struct {
	loader           abstraction.UrlLoader
	extractor        abstraction.ConfigsExtractor
	writer           abstraction.ConfigsWriter
	configRepository repository.ConfigRepository
	groupRepository  repository.GroupRepository
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

	var formattedConfigs []string
	for _, rawConfig := range rawConfigs {
		configType := strings.Split(rawConfig, "://")[0]
		formatter, err := creator.formatterFactory.Fetch(configType)
		if err != nil {
			continue
		}

		formattedConfig, err := formatter.Format(rawConfig)
		if err == nil {
			formattedConfigs = append(formattedConfigs, formattedConfig)
		}
	}
	if len(formattedConfigs) == 0 {
		errText := creator.configuration.Errors.ConfigFormatError
		return errors.New(errText)
	}

	goroutinesMaxCount := creator.configuration.GoroutinesMax
	semaphore := make(chan struct{}, goroutinesMaxCount)
	var wg sync.WaitGroup

	newGroupId, err := creator.groupRepository.CreateGroup(groupName)
	if err != nil {
		errText := creator.configuration.Errors.GroupCreatingError
		return fmt.Errorf("%v %w", errText, err)
	}

	v2ConfigsPath := creator.configuration.V2.ConfigurationPath
	groupPath := path.Join(v2ConfigsPath, groupName)
	err = os.Mkdir(groupPath, 0755)
	if err != nil {
		errText := creator.configuration.Errors.DirectoryCreatingError
		return fmt.Errorf("%v %w", errText, err)
	}

	for i, formattedConfig := range formattedConfigs {

		wg.Add(1)
		go func(config string, index int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			configPath := path.Join(groupPath, strconv.Itoa(index))
			err = creator.writer.Write(config, configPath)
			if err == nil {
				creator.configRepository.CreateConfig(newGroupId, configPath)
			}
		}(formattedConfig, i)
	}
	wg.Wait()

	return nil
}
