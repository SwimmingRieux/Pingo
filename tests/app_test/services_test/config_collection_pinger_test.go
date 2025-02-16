package services_test

import (
	"github.com/stretchr/testify/mock"
	"net"
	configurations "pingo/configs"
	"pingo/internal/app/services"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/structs"
	"sync"
	"testing"
)

type MockConfigPinger struct {
	mock.Mock
}

func (m *MockConfigPinger) Ping(config entities.Config, domain structs.DomainWithRank, listener net.Listener, configScoresMap *sync.Map) {
	m.Called(config, domain, listener, configScoresMap)
}

var realConfig, _ = configurations.NewConfig()

var pingAllConfigTests = []pingAllConfigTest{
	{
		name:         "should call ConfigPinger.Ping() exactly once for each pair of (config,domain) when parameters are valid",
		configs:      make([]entities.Config, 4),
		domains:      make([]structs.DomainWithRank, 5),
		listeners:    make([]net.Listener, 4),
		testFunction: PingAllConfigsCallEachPairOnceWhenParametersAreValid,
	},
	{
		name:         "should not call ConfigPinger.Ping() at all when configs count is 0",
		configs:      []entities.Config{},
		domains:      make([]structs.DomainWithRank, 5),
		listeners:    []net.Listener{},
		testFunction: PingAllConfigsCallEachPairOnceWhenParametersAreValid,
	},
	{
		name:         "should not call ConfigPinger.Ping() at all when domains count is 0",
		configs:      make([]entities.Config, 3),
		domains:      []structs.DomainWithRank{},
		listeners:    make([]net.Listener, 3),
		testFunction: PingAllConfigsCallEachPairOnceWhenParametersAreValid,
	},
}

func TestPingAllConfigs(t *testing.T) {
	t.Parallel()
	for _, testCase := range pingAllConfigTests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.testFunction(t, testCase.configs, testCase.domains, testCase.listeners)
		})
	}
}

func PingAllConfigsCallEachPairOnceWhenParametersAreValid(t *testing.T, configs []entities.Config, domains []structs.DomainWithRank, listeners []net.Listener) {
	// Arrange
	mockConfigPinger := new(MockConfigPinger)
	for _, config := range configs {
		for _, domain := range domains {
			mockConfigPinger.On("Ping", config, domain, mock.Anything, mock.Anything).Once()
		}
	}

	service := services.NewConfigCollectionPinger(mockConfigPinger, *realConfig)

	configScoresMap := &sync.Map{}
	var wg sync.WaitGroup
	// Act
	service.PingAllConfigs(configs, domains, &wg, listeners, configScoresMap)
	wg.Wait()

	// Assert
	mockConfigPinger.AssertExpectations(t)
}

type pingAllConfigTest struct {
	name         string
	configs      []entities.Config
	domains      []structs.DomainWithRank
	listeners    []net.Listener
	testFunction func(t *testing.T, configs []entities.Config, domains []structs.DomainWithRank, listeners []net.Listener)
}
