package app

import (
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

func (this *ConfigCreator) Create(input string) error {
	var configsString string
	var err error

	if strings.HasPrefix(input, "http") && !strings.Contains(input, " ") {
		configsString, err = this.Loader.Load(input)
		if err != nil {
			return fmt.Errorf("Could not load from the link: ", err)
		}
	} else {
		configsString = input
	}

	groupName, rawConfigs := this.Extractor.Extract(configsString)
	if len(rawConfigs) == 0 {
		return fmt.Errorf("Could not found any config")
	}

	var formattedConfigs []string
	for _, rawConfig := range rawConfigs {
		formattedConfig, err := this.Formatter.Format(rawConfig)
		if err == nil {
			formattedConfigs = append(formattedConfigs, formattedConfig)
		}
	}
	if len(formattedConfigs) == 0 {
		return fmt.Errorf("Could not format any config")
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)

	newGroupId, err := this.Repository.CreateGroup(groupName)
	if err != nil {
		return fmt.Errorf("Error creating grpup:", err)
	}
	groupPath := path.Join("default_path", groupName)
	err := os.Mkdir(groupPath, 0755)
	if err != nil {
		return fmt.Errorf("Error creating directory:", err)
	}

	for i, formattedConfig := range formattedConfigs {

		wg.Add(1)
		go func(config string, index int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			configPath := path.Join(groupPath, strconv.Itoa(index))
			err = this.Writer.Write(config, configPath)
			if err == nil {
				this.Repository.CreateConfig(newGroupId, configPath)
			}
		}(formattedConfig, i)
	}
	wg.Wait()

	return nil
}
