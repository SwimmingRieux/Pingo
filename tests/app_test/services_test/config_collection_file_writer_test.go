package services_test

import (
	"path"
	"pingo/internal/app/services"
	"pingo/internal/domain/structs"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockConfigFileWriter struct {
	mock.Mock
}

func (m *MockConfigFileWriter) Write(config string, filePath string) error {
	args := m.Called(config, filePath)
	return args.Error(0)
}

type MockRepositoryConfigCreator struct {
	mock.Mock
}

func (m *MockRepositoryConfigCreator) CreateConfig(groupId int, path string, configType string) (int, error) {
	args := m.Called(groupId, path, configType)
	return 0, args.Error(0)
}

func TestWriteConfigsToFiles(t *testing.T) {
	// Arrange
	mockWriter := new(MockConfigFileWriter)
	mockRepo := new(MockRepositoryConfigCreator)

	formattedConfigs := []structs.FormattedConfigAndType{
		{FormattedConfig: "config1", Type: "json"},
		{FormattedConfig: "config2", Type: "yaml"},
		{FormattedConfig: "config3", Type: "xml"},
	}
	groupPath := "/test/path"
	newGroupId := 123
	for i, config := range formattedConfigs {
		configPath := path.Join(groupPath, strconv.Itoa(i))
		mockWriter.On("Write", config.FormattedConfig, configPath).Return(nil).Once()
		mockRepo.On("CreateConfig", newGroupId, configPath, config.Type).Return(nil).Once()
	}
	writer := services.NewConfigCollectionFileWriter(mockWriter, mockRepo, ConfigForTest)
	var wg sync.WaitGroup

	// Act
	writer.WriteConfigsToFiles(formattedConfigs, &wg, groupPath, newGroupId)
	wg.Wait()

	// Assert
	mockWriter.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
