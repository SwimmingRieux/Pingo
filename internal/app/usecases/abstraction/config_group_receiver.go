package abstraction

import "pingo/internal/domain/entities"

type ConfigGroupReceiver interface {
	Get(id int) (entities.Config, error)
}
