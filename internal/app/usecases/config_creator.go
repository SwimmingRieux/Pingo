package app

import (
	"errors"
	"fmt"
	"os"
	"path"
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
	formatter        abstraction.ConfigsFormatter
	configRepository repository.ConfigRepository
	groupRepository  repository.GroupRepository
}

func (creator *ConfigCreator) Create(input string) error {
	var configsString string
	var err error

	if strings.HasPrefix(input, "http") && !strings.Contains(input, " ") {
		configsString, err = creator.loader.Load(input)
		if err != nil {
			return fmt.Errorf("could not load from the link: %w", err)
		}
	} else {
		configsString = input
	}

	groupName, rawConfigs := creator.extractor.Extract(configsString)
	if len(rawConfigs) == 0 {
		return errors.New("could not found any config")
	}

	var formattedConfigs []string
	for _, rawConfig := range rawConfigs {
		formattedConfig, err := creator.formatter.Format(rawConfig)
		if err == nil {
			formattedConfigs = append(formattedConfigs, formattedConfig)
		}
	}
	if len(formattedConfigs) == 0 {
		return errors.New("could not format any config")
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)

	newGroupId, err := creator.groupRepository.CreateGroup(groupName)
	if err != nil {
		return fmt.Errorf("error creating grpup: %w", err)
	}
	groupPath := path.Join("default_path", groupName)
	err = os.Mkdir(groupPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
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
