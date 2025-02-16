package services_test

import (
	"github.com/stretchr/testify/mock"
	"pingo/internal/app/services"
	"pingo/internal/domain/dtos"
	"pingo/internal/domain/entities"
	"sync"
	"testing"
)

type MockConfigRepository struct {
	mock.Mock
}

func (m *MockConfigRepository) UpdateConfig(id int, configDto dtos.UpdateConfigDto) {
	m.Called(id, configDto)
}

type scoreWriterTest struct {
	name              string
	configs           []entities.Config
	configScoresMap   *sync.Map
	expectedConfigDto []dtos.UpdateConfigDto
}

var testCases = []scoreWriterTest{
	{
		name:              "should not update any configurations when no configurations are provided",
		configs:           []entities.Config{},
		configScoresMap:   &sync.Map{},
		expectedConfigDto: []dtos.UpdateConfigDto{},
	},
	{
		name: "should update a single configuration with a valid score when one configuration is provided and its score exists in the configScoresMap",
		configs: []entities.Config{
			{ConfigId: 1, Type: "type1", Path: "path1", Score: 0.0},
		},
		configScoresMap: func() *sync.Map {
			m := &sync.Map{}
			m.Store(1, 0.8) // Using ConfigId as the key
			return m
		}(),
		expectedConfigDto: []dtos.UpdateConfigDto{
			{Type: "type1", Path: "path1", Score: 0.8},
		},
	},
	{
		name: "should update multiple configurations with valid scores when multiple configurations are provided and their scores exist in the configScoresMap",
		configs: []entities.Config{
			{ConfigId: 1, Type: "type1", Path: "path1", Score: 0.0},
			{ConfigId: 2, Type: "type2", Path: "path2", Score: 0.0},
		},
		configScoresMap: func() *sync.Map {
			m := &sync.Map{}
			m.Store(1, 0.8) // Using ConfigId as the key
			m.Store(2, 0.9) // Using ConfigId as the key
			return m
		}(),
		expectedConfigDto: []dtos.UpdateConfigDto{
			{Type: "type1", Path: "path1", Score: 0.8},
			{Type: "type2", Path: "path2", Score: 0.9},
		},
	},
	{
		name: "should only update configurations with scores available in the configScoresMap when some configurations have missing scores",
		configs: []entities.Config{
			{ConfigId: 1, Type: "type1", Path: "path1", Score: 0.0},
			{ConfigId: 2, Type: "type2", Path: "path2", Score: 0.0},
		},
		configScoresMap: func() *sync.Map {
			m := &sync.Map{}
			m.Store(1, 0.8) // Using ConfigId as the key
			return m
		}(),
		expectedConfigDto: []dtos.UpdateConfigDto{
			{Type: "type1", Path: "path1", Score: 0.8},
		},
	},
}

func TestWriteScoresToDb(t *testing.T) {
	t.Parallel()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockConfigRepository := new(MockConfigRepository)
			service := services.NewConfigScoreWriter(mockConfigRepository)
			for i, expectedConfigItem := range testCase.expectedConfigDto {
				mockConfigRepository.On("UpdateConfig", testCase.configs[i].ConfigId, expectedConfigItem).Once()
			}

			// Act
			service.WriteScoresToDb(testCase.configs, testCase.configScoresMap)

			// Assert
			mockConfigRepository.AssertExpectations(t)
		})
	}
}
