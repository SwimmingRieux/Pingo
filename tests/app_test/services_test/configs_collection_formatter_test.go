package services_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pingo/internal/app/services"
	"pingo/internal/app/services/abstraction"
	"testing"
)

type configsCollectionFormatterTest struct {
	name                   string
	canBeFormatted         bool
	expectedContainedError error
}

type MockFormatterFactory struct {
	mock.Mock
}

func (m *MockFormatterFactory) Fetch(formatterType string) (abstraction.ConfigsFormatter, error) {
	args := m.Called(formatterType)
	return args.Get(0).(abstraction.ConfigsFormatter), args.Error(1)
}

type MockFormatter struct {
	mock.Mock
}

func (m *MockFormatter) Format(rawConfig string) (string, error) {
	args := m.Called(rawConfig)
	return args.String(0), args.Error(1)
}

var configsCollectionFormatterTests = []configsCollectionFormatterTest{
	{
		name:                   "should return valid slice and call formatter and return nil error when formatter returns valid formatted configs",
		canBeFormatted:         true,
		expectedContainedError: nil,
	},
	{
		name:                   "should return error and empty slice and call formatter and return nil error when formatter returns error for all the raw configs",
		canBeFormatted:         false,
		expectedContainedError: fmt.Errorf("%v", ConfigForTest.Errors.ConfigFormatError),
	},
}

func TestFormatCollection(t *testing.T) {
	t.Parallel()
	for _, testCase := range configsCollectionFormatterTests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			var rawConfigs = []string{"1", "2", "3", "4", "5"}

			mockFormatter := new(MockFormatter)
			if testCase.canBeFormatted {
				mockFormatter.On("Format", mock.Anything).Return("validString", nil)
			} else {
				mockFormatter.On("Format", mock.Anything).Return("", errors.New(""))
			}
			mockFormatterFactory := new(MockFormatterFactory)
			mockFormatterFactory.On("Fetch", mock.Anything).Return(mockFormatter, nil).Times(len(rawConfigs))
			service := services.NewConfigsCollectionFormatter(mockFormatterFactory, ConfigForTest)

			// Act
			_, err := service.FormatCollection(rawConfigs)

			// Assert
			assert.Equal(t, testCase.expectedContainedError, err)
			mockFormatter.AssertExpectations(t)
			mockFormatterFactory.AssertExpectations(t)
		})
	}
}
