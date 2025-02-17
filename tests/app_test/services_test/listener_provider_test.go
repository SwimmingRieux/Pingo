package services_test

import (
	"github.com/stretchr/testify/assert"
	"pingo/internal/app/services"
	"testing"
)

type listenerProviderTest struct {
	name          string
	configsLength int
	testFunction  func(t *testing.T, configsLength int)
}

var listenerProviderTests = []listenerProviderTest{
	{
		name:          "should return an array of listeners when the length is not bigger than iteration limit",
		configsLength: 5,
		testFunction:  listenerProviderSuccessfulTest,
	},
	{
		name:          "should return error when length is bigger than iteration limit",
		configsLength: 1000000,
		testFunction:  listenerProviderFailedTest,
	},
}

func TestListenerProvider(t *testing.T) {
	t.Parallel()
	for _, testCase := range listenerProviderTests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			testCase.testFunction(t, testCase.configsLength)
		})
	}
}

func listenerProviderSuccessfulTest(t *testing.T, configsLength int) {
	// Arrange
	provider := services.NewListenerProvider(ConfigForTest)
	// Act
	listeners, err := provider.GetListeners(configsLength)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, listeners, configsLength)

	for _, listener := range listeners {
		err := listener.Close()
		assert.NoError(t, err)
	}
}

func listenerProviderFailedTest(t *testing.T, configsLength int) {
	// Arrange
	provider := services.NewListenerProvider(ConfigForTest)
	// Act
	listeners, err := provider.GetListeners(configsLength)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, listeners)
	assert.EqualError(t, err, ConfigForTest.Errors.NotEnoughPortsFound)
}
