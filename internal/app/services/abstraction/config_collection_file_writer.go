package abstraction

import (
	"pingo/internal/domain/structs"
	"sync"
)

type ConfigCollectionFileWriter interface {
	WriteConfigsToFiles(formattedConfigs []structs.FormattedConfigAndType, wg *sync.WaitGroup, groupPath string, newGroupId int)
}
