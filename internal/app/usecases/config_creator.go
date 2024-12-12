package usecases

import (
	"errors"
	"fmt"
	"os"
	"path"
	configsAbstraction "pingo/configs/abstraction"
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
	configReader     configsAbstraction.Config
}

func (creator *ConfigCreator) Create(input string) error {
	var configsString string
	var err error

	if strings.HasPrefix(input, "http") && !strings.Contains(input, " ") {
		configsString, err = creator.loader.Load(input)
		if err != nil {
			errText, _ := creator.configReader.Get("errors.load_from_link_error")
			return fmt.Errorf("%v %w", errText, err)
		}
	} else {
		configsString = input
	}

	groupName, rawConfigs := creator.extractor.Extract(configsString)
	if len(rawConfigs) == 0 {
		errText, _ := creator.configReader.Get("errors.config_not_found")
		return errors.New(errText)
	}

	var formattedConfigs []string
	for _, rawConfig := range rawConfigs {
		formattedConfig, err := creator.formatter.Format(rawConfig)
		if err == nil {
			formattedConfigs = append(formattedConfigs, formattedConfig)
		}
	}
	if len(formattedConfigs) == 0 {
		errText, _ := creator.configReader.Get("errors.config_format_error")
		return errors.New(errText)
	}

	stringCount, _ := creator.configReader.Get("goroutines_max")
	goroutinesMaxCount, _ := strconv.Atoi(stringCount)
	semaphore := make(chan struct{}, goroutinesMaxCount)
	var wg sync.WaitGroup

	newGroupId, err := creator.groupRepository.CreateGroup(groupName)
	if err != nil {
		errText, _ := creator.configReader.Get("errors.group_creating_error")
		return fmt.Errorf("%v %w", errText, err)
	}

	v2ConfigsPath, _ := creator.configReader.Get("v2.config_path")
	groupPath := path.Join(v2ConfigsPath, groupName)
	err = os.Mkdir(groupPath, 0755)
	if err != nil {
		errText, _ := creator.configReader.Get("errors.directory_creating_error")
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
