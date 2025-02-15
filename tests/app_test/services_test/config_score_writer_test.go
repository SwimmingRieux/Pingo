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
	name            string
	configs         []entities.Config
	configScoresMap *sync.Map
}

var testCases = []scoreWriterTest{
	{
		name:            "",
		configs:         nil,
		configScoresMap: nil,
	},
}

func TestWriteScoresToDb(t *testing.T) {
	for _, testCase := range testCases {
		t.Parallel()
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockConfigRepository := new(MockConfigRepository)
			service := services.NewConfigScoreWriter(mockConfigRepository)
		})
	}
}
