package services

import (
	"pingo/internal/domain/dtos"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/repository"
	"sync"
)

type ConfigScoreWriter struct {
	configRepository repository.ConfigRepository
}

func (s *ConfigScoreWriter) WriteScoresToDb(configs []entities.Config, configScoresMap *sync.Map) {
	for _, config := range configs {
		value, ok := configScoresMap.Load(config)
		if ok {
			configDto := dtos.UpdateConfigDto{Type: config.Type, Path: config.Path, Score: value.(float64)}
			s.configRepository.UpdateConfig(config.ConfigId, configDto)
		} else {
			configDto := dtos.UpdateConfigDto{Type: config.Type, Path: config.Path, Score: 0}
			s.configRepository.UpdateConfig(config.ConfigId, configDto)
		}
	}
}
