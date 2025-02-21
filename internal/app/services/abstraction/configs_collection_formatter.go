package abstraction

import "pingo/internal/domain/structs"

type ConfigsCollectionFormatter interface {
	FormatCollection(rawConfigs []string) ([]structs.FormattedConfigAndType, error)
}
