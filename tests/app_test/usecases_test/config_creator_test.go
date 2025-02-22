package usecases_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/app/usecases"
	"pingo/internal/domain/repository"
	"pingo/internal/domain/structs"
	"sync"
	"testing"
)

type mockUrlLoader struct {
	abstraction.UrlLoader
	loadFunc func(string) (string, error)
}

func (m *mockUrlLoader) Load(url string) (string, error) {
	return m.loadFunc(url)
}

type mockConfigsExtractor struct {
	abstraction.ConfigsExtractor
	extractFunc func(string) (string, []string)
}

func (m *mockConfigsExtractor) Extract(input string) (string, []string) {
	return m.extractFunc(input)
}

type mockConfigCollectionFileWriter struct {
	abstraction.ConfigCollectionFileWriter
	writeConfigsToFilesFunc func([]structs.FormattedConfigAndType, *sync.WaitGroup, string, int)
}

func (m *mockConfigCollectionFileWriter) WriteConfigsToFiles(configs []structs.FormattedConfigAndType, wg *sync.WaitGroup, path string, groupId int) {
	m.writeConfigsToFilesFunc(configs, wg, path, groupId)
}

type mockRepositoryGroupCreator struct {
	repository.RepositoryGroupCreator
	createGroupFunc func(string) (int, error)
}

func (m *mockRepositoryGroupCreator) CreateGroup(name string) (int, error) {
	return m.createGroupFunc(name)
}

type mockConfigsCollectionFormatter struct {
	abstraction.ConfigsCollectionFormatter
	formatCollectionFunc func([]string) ([]structs.FormattedConfigAndType, error)
}

func (m *mockConfigsCollectionFormatter) FormatCollection(configs []string) ([]structs.FormattedConfigAndType, error) {
	return m.formatCollectionFunc(configs)
}

type configCreatorFailedTest struct {
	name                         string
	loader                       abstraction.UrlLoader
	extractor                    abstraction.ConfigsExtractor
	collectionWriter             abstraction.ConfigCollectionFileWriter
	groupRepository              repository.RepositoryGroupCreator
	collectionFormatter          abstraction.ConfigsCollectionFormatter
	expectedContainedErrorString string
}

type configCreatorSuccessTest struct {
	name                string
	inputString         string
	loader              abstraction.UrlLoader
	extractor           abstraction.ConfigsExtractor
	collectionWriter    abstraction.ConfigCollectionFileWriter
	groupRepository     repository.RepositoryGroupCreator
	collectionFormatter abstraction.ConfigsCollectionFormatter
}

var configCreatorFailedTests = []configCreatorFailedTest{
	{
		name: "should return error when url loader fails",
		loader: &mockUrlLoader{
			loadFunc: func(url string) (string, error) {
				return "", errors.New("load error")
			},
		},
		extractor: &mockConfigsExtractor{
			extractFunc: func(input string) (string, []string) {
				return "group", []string{"config1", "config2"}
			},
		},
		collectionWriter: &mockConfigCollectionFileWriter{
			writeConfigsToFilesFunc: func(configs []structs.FormattedConfigAndType, wg *sync.WaitGroup, path string, groupId int) {
			},
		},
		groupRepository: &mockRepositoryGroupCreator{
			createGroupFunc: func(name string) (int, error) {
				return 1, nil
			},
		},
		collectionFormatter: &mockConfigsCollectionFormatter{
			formatCollectionFunc: func(configs []string) ([]structs.FormattedConfigAndType, error) {
				return []structs.FormattedConfigAndType{}, nil
			},
		},
		expectedContainedErrorString: "load error",
	},
	{
		name: "should return error when configs not found",
		loader: &mockUrlLoader{
			loadFunc: func(url string) (string, error) {
				return "configs", nil
			},
		},
		extractor: &mockConfigsExtractor{
			extractFunc: func(input string) (string, []string) {
				return "group", []string{}
			},
		},
		collectionWriter: &mockConfigCollectionFileWriter{
			writeConfigsToFilesFunc: func(configs []structs.FormattedConfigAndType, wg *sync.WaitGroup, path string, groupId int) {
			},
		},
		groupRepository: &mockRepositoryGroupCreator{
			createGroupFunc: func(name string) (int, error) {
				return 1, nil
			},
		},
		collectionFormatter: &mockConfigsCollectionFormatter{
			formatCollectionFunc: func(configs []string) ([]structs.FormattedConfigAndType, error) {
				return []structs.FormattedConfigAndType{}, nil
			},
		},
		expectedContainedErrorString: "config not found",
	},
}

var configCreatorSuccessTests = []configCreatorSuccessTest{
	{
		name: "should create configs when input is valid url",
		loader: &mockUrlLoader{
			loadFunc: func(url string) (string, error) {
				return "configs", nil
			},
		},
		extractor: &mockConfigsExtractor{
			extractFunc: func(input string) (string, []string) {
				return "group", []string{"config1", "config2"}
			},
		},
		collectionWriter: &mockConfigCollectionFileWriter{
			writeConfigsToFilesFunc: func(configs []structs.FormattedConfigAndType, wg *sync.WaitGroup, path string, groupId int) {
			},
		},
		groupRepository: &mockRepositoryGroupCreator{
			createGroupFunc: func(name string) (int, error) {
				return 1, nil
			},
		},
		collectionFormatter: &mockConfigsCollectionFormatter{
			formatCollectionFunc: func(configs []string) ([]structs.FormattedConfigAndType, error) {
				return []structs.FormattedConfigAndType{{FormattedConfig: "config1"}, {FormattedConfig: "config2"}}, nil
			},
		},
	},
	{
		name: "should create configs when input is valid config string",
		loader: &mockUrlLoader{
			loadFunc: func(url string) (string, error) {
				return "", nil
			},
		},
		extractor: &mockConfigsExtractor{
			extractFunc: func(input string) (string, []string) {
				return "group", []string{"config1", "config2"}
			},
		},
		collectionWriter: &mockConfigCollectionFileWriter{
			writeConfigsToFilesFunc: func(configs []structs.FormattedConfigAndType, wg *sync.WaitGroup, path string, groupId int) {
			},
		},
		groupRepository: &mockRepositoryGroupCreator{
			createGroupFunc: func(name string) (int, error) {
				return 1, nil
			},
		},
		collectionFormatter: &mockConfigsCollectionFormatter{
			formatCollectionFunc: func(configs []string) ([]structs.FormattedConfigAndType, error) {
				return []structs.FormattedConfigAndType{{FormattedConfig: "config1"}, {FormattedConfig: "config2"}}, nil
			},
		},
	},
}

func testFailedCreate(t *testing.T, test configCreatorFailedTest) {
	// Arrange
	defer func() {
		pingoPath := os.Getenv("PINGO_PATH")
		v2ConfigsPath := ConfigForTest.V2.ConfigurationPath
		groupPath := path.Join(pingoPath, v2ConfigsPath, "group")
		_ = os.RemoveAll(groupPath)
	}()
	creator := usecases.NewConfigCreator(test.loader, test.extractor, test.collectionWriter, test.groupRepository, ConfigForTest, test.collectionFormatter)
	// Act
	err := creator.Create("http://example.com")
	// Assert
	if err == nil {
		t.Fatalf("expected error, didn't get any")
	}
	assert.Contains(t, err.Error(), test.expectedContainedErrorString)
}

func testSuccessCreate(t *testing.T, test configCreatorSuccessTest) {
	// Arrange
	defer func() {
		pingoPath := os.Getenv("PINGO_PATH")
		v2ConfigsPath := ConfigForTest.V2.ConfigurationPath
		groupPath := path.Join(pingoPath, v2ConfigsPath, "group")
		_ = os.RemoveAll(groupPath)
	}()
	creator := usecases.NewConfigCreator(test.loader, test.extractor, test.collectionWriter, test.groupRepository, ConfigForTest, test.collectionFormatter)
	// Act
	err := creator.Create(test.inputString)
	// Assert
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}

}

func TestCreate(t *testing.T) {
	t.Parallel()
	for _, testCase := range configCreatorFailedTests {
		t.Run(testCase.name, func(t *testing.T) {
			testFailedCreate(t, testCase)
		})
	}
	for _, testCase := range configCreatorSuccessTests {
		t.Run(testCase.name, func(t *testing.T) {
			testSuccessCreate(t, testCase)
		})
	}
}
