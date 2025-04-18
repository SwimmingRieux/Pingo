package services_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"pingo/internal/app/services"
	"testing"
)

type configsFileWriterTest struct {
	name                       string
	jsonConfig                 string
	path                       string
	expectedContainedErrorText string
	testFunction               func(t *testing.T, jsonConfig string, path string, decoded string)
}

var configsFileWriterTests = []configsFileWriterTest{
	{
		name:                       "should create file and write content successfully when valid JSON and path are provided",
		jsonConfig:                 `{"key": "value"}`,
		path:                       "test_config.json",
		expectedContainedErrorText: "",
		testFunction:               ConfigsFileWriterSuccessfulTest,
	},
	{
		name:                       "should return error when path is empty",
		jsonConfig:                 `{"key": "value"}`,
		path:                       "",
		expectedContainedErrorText: ConfigForTest.Errors.FileCreatingError,
		testFunction:               ConfigFileWriterFailedTest,
	},
	{
		name:                       "should return error when path is invalid (e.g., directory does not exist)",
		jsonConfig:                 `{"key": "value"}`,
		path:                       "/nonexistent_directory_nsdjfAFBNKDaaefb/test_config.json",
		expectedContainedErrorText: ConfigForTest.Errors.FileCreatingError,
		testFunction:               ConfigFileWriterFailedTest,
	},
}

func TestConfigsFileWriter(t *testing.T) {
	t.Parallel()
	for _, testCase := range configsFileWriterTests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			testCase.testFunction(t, testCase.jsonConfig, testCase.path, testCase.expectedContainedErrorText)
		})
	}
}

func ConfigFileWriterFailedTest(t *testing.T, jsonConfig string, path string, expectedContainedErrorText string) {
	// Arrange
	service := services.NewConfigFileWriter(ConfigForTest)
	// Act
	err := service.Write(jsonConfig, path)
	// Assert
	assert.ErrorContains(t, err, expectedContainedErrorText)
}

func ConfigsFileWriterSuccessfulTest(t *testing.T, jsonConfig string, path string, expectedContainedErrorText string) {
	// Arrange
	service := services.NewConfigFileWriter(ConfigForTest)
	// Act
	err := service.Write(jsonConfig, path)
	// Assert
	assert.NoError(t, err)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}
	if string(content) != jsonConfig {
		t.Errorf("File content does not match. Expected: %s, Got: %s", jsonConfig, string(content))
	}
	err = os.Remove(path)
	if err != nil {
		t.Errorf("Failed to clean up test file: %v", err)
	}
}
