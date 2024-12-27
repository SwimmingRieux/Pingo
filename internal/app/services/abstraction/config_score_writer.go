package abstraction

import (
	"pingo/internal/domain/entities"
	"sync"
)

type ConfigScoreWriter interface {
	WriteScoresToDb(configs []entities.Config, configScoresMap *sync.Map)
}
