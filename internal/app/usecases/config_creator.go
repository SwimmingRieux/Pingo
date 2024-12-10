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
	Loader     abstraction.UrlLoader
	Extractor  abstraction.ConfigsExtractor
	Writer     abstraction.ConfigsWriter
	Formatter  abstraction.ConfigsFormatter
	Repository repository.ConfigRepository
}

func (creator *ConfigCreator) Create(input string) error {
	var configsString string
	var err error

	if strings.HasPrefix(input, "http") && !strings.Contains(input, " ") {
		configsString, err = creator.Loader.Load(input)
		if err != nil {
			return fmt.Errorf("could not load from the link: %w", err)
		}
	} else {
		configsString = input
	}

	groupName, rawConfigs := creator.Extractor.Extract(configsString)
	if len(rawConfigs) == 0 {
		return errors.New("could not found any config")
	}

	var formattedConfigs []string
	for _, rawConfig := range rawConfigs {
		formattedConfig, err := creator.Formatter.Format(rawConfig)
		if err == nil {
			formattedConfigs = append(formattedConfigs, formattedConfig)
		}
	}
	if len(formattedConfigs) == 0 {
		return errors.New("could not format any config")
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)

	newGroupId, err := creator.Repository.CreateGroup(groupName)
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
			err = creator.Writer.Write(config, configPath)
			if err == nil {
				creator.Repository.CreateConfig(newGroupId, configPath)
			}
		}(formattedConfig, i)
	}
	wg.Wait()

	return nil
}
