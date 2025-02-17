package services

import (
	"pingo/internal/domain/dtos"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/repository"
	"sync"
)

type ConfigScoreWriter struct {
	configRepository repository.RepositoryConfigUpdater
}

func NewConfigScoreWriter(configRepository repository.RepositoryConfigUpdater) *ConfigScoreWriter {
	return &ConfigScoreWriter{
		configRepository: configRepository,
	}
}

func (s *ConfigScoreWriter) WriteScoresToDb(configs []entities.Config, configScoresMap *sync.Map) {
	for _, config := range configs {
		value, ok := configScoresMap.Load(config.ConfigId)
		if ok {
			score := value.(float64)
			configDto := dtos.UpdateConfigDto{Type: config.Type, Path: config.Path, Score: score}
			s.configRepository.UpdateConfig(config.ConfigId, configDto)
		}
	}
}
