package services_test

import (
	"github.com/stretchr/testify/assert"
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

func TestPingAllConfigs(t *testing.T) {
	mockConfigPinger := new(MockConfigPinger)
	mockConfigPinger.On("Ping", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	realConfig, err := configurations.NewConfig()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	service := services.NewConfigCollectionPinger(mockConfigPinger, *realConfig)

	configs := []entities.Config{
		{}, {},
	}
	domains := []structs.DomainWithRank{
		{},
		{},
	}
	var listeners []net.Listener
	for range configs {
		listeners = append(listeners, nil)
	}
	configScoresMap := &sync.Map{}
	var wg sync.WaitGroup
	service.PingAllConfigs(configs, domains, &wg, listeners, configScoresMap)

	wg.Wait()

	for _, config := range configs {
		for _, domain := range domains {
			mockConfigPinger.AssertCalled(t, "Ping", config, domain, mock.Anything, configScoresMap)
		}
	}

	assert.Equal(t, len(configs)*len(domains), len(mockConfigPinger.Calls))
}
