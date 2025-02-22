package usecases_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"path/filepath"
	"pingo/internal/app/usecases"
	"pingo/internal/domain/entities"
	"testing"
)

type MockDeleter struct {
	mock.Mock
}

func (m *MockDeleter) DeleteConfig(configID int) error {
	args := m.Called(configID)
	return args.Error(0)
}

type MockRetriever struct {
	mock.Mock
}

func (m *MockRetriever) GetConfig(configID int) (entities.Config, error) {
	args := m.Called(configID)
	return args.Get(0).(entities.Config), args.Error(1)
}

var configRemoverTests = []struct {
	name         string
	testFunction func(t *testing.T)
}{
	{
		name:         "should return error when can not delete config from db",
		testFunction: testDeleteConfigError,
	},
	{
		name:         "should return error when the id is invalid",
		testFunction: testInvalidIDError,
	},
	{
		name:         "should return error when can not remove the config file",
		testFunction: testFileRemoveError,
	},
	{
		name:         "should remove successfully when no errors occur",
		testFunction: testRemoveSuccess,
	},
}

func testDeleteConfigError(t *testing.T) {
	t.Parallel()
	// Arrange
	mockDeleter := new(MockDeleter)
	mockRetriever := new(MockRetriever)
	config := entities.Config{Path: "testpath"}
	mockRetriever.On("GetConfig", 1).Return(config, nil)
	mockDeleter.On("DeleteConfig", 1).Return(errors.New("db error"))
	remover := usecases.NewConfigRemover(mockDeleter, mockRetriever, ConfigForTest)

	// Act
	err := remover.Remove(1)

	// Assert
	assert.ErrorContains(t, err, ConfigForTest.Errors.ConfigRemoveError)
}

func testInvalidIDError(t *testing.T) {
	t.Parallel()
	// Arrange
	mockRetriever := new(MockRetriever)
	mockRetriever.On("GetConfig", 2).Return(entities.Config{}, errors.New("not found"))
	remover := usecases.NewConfigRemover(new(MockDeleter), mockRetriever, ConfigForTest)

	// Act
	err := remover.Remove(2)

	// Assert
	assert.ErrorContains(t, err, ConfigForTest.Errors.ConfigNotFound)
}

func testFileRemoveError(t *testing.T) {
	t.Parallel()
	// Arrange
	mockDeleter := new(MockDeleter)
	mockRetriever := new(MockRetriever)
	config := entities.Config{Path: "apaththatdoesntexist"}
	mockRetriever.On("GetConfig", 3).Return(config, nil)
	mockDeleter.On("DeleteConfig", 3).Return(nil)
	remover := usecases.NewConfigRemover(mockDeleter, mockRetriever, ConfigForTest)
	// Act
	err := remover.Remove(3)

	// Assert
	assert.ErrorContains(t, err, ConfigForTest.Errors.FileRemoveError)
}

func testRemoveSuccess(t *testing.T) {
	t.Parallel()
	// Arrange
	mockDeleter := new(MockDeleter)
	mockRetriever := new(MockRetriever)
	config := entities.Config{Path: "temporarypathforafilethatwillbedeleted.txt"}
	mockRetriever.On("GetConfig", 4).Return(config, nil)
	mockDeleter.On("DeleteConfig", 4).Return(nil)

	pingoPath := os.Getenv("PINGO_PATH")
	filePath := filepath.Join(pingoPath, ConfigForTest.V2.ConfigurationPath, config.Path)
	tempFile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer func(name string) {
		os.Remove(name)
	}(tempFile.Name())

	remover := usecases.NewConfigRemover(mockDeleter, mockRetriever, ConfigForTest)

	// Act
	err = remover.Remove(4)

	// Assert
	assert.NoError(t, err)
}

func TestRemove(t *testing.T) {
	t.Parallel()
	for _, testCase := range configRemoverTests {
		t.Run(testCase.name, testCase.testFunction)
	}
}
