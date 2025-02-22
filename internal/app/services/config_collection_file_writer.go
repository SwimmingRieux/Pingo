package services

import (
	"path"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/repository"
	"pingo/internal/domain/structs"
	"strconv"
	"sync"
)

type ConfigCollectionFileWriter struct {
	singleFileWriter abstraction.ConfigFileWriter
	configRepository repository.RepositoryConfigCreator
	configuration    *configs.Configuration
}

func NewConfigCollectionFileWriter(singleFileWriter abstraction.ConfigFileWriter, configRepository repository.RepositoryConfigCreator, configuration *configs.Configuration) *ConfigCollectionFileWriter {
	return &ConfigCollectionFileWriter{
		singleFileWriter: singleFileWriter,
		configRepository: configRepository,
		configuration:    configuration,
	}
}

func (collectionWriter *ConfigCollectionFileWriter) WriteConfigsToFiles(formattedConfigs []structs.FormattedConfigAndType, wg *sync.WaitGroup,
	groupPath string, newGroupId int) {

	goroutinesMaxCount := collectionWriter.configuration.GoroutinesMax
	semaphore := make(chan struct{}, goroutinesMaxCount)

	for i, formattedConfig := range formattedConfigs {
		wg.Add(1)
		go func(i int, formattedConfig structs.FormattedConfigAndType) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			configPath := path.Join(groupPath, strconv.Itoa(i))
			err := collectionWriter.singleFileWriter.Write(formattedConfig.FormattedConfig, configPath)
			if err == nil {
				collectionWriter.configRepository.CreateConfig(newGroupId, configPath, formattedConfig.Type)
			}
		}(i, formattedConfig)
	}
}
