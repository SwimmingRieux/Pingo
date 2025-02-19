package usecases_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/app/usecases"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/structs"
	"sync"
	"testing"
)

type configOrganizerTest struct {
	name              string
	domainsCountLimit int
	configs           []entities.Config
	domainsWithRank   []structs.DomainWithRank
}

type MockNetListener struct{}

func (dl MockNetListener) Accept() (net.Conn, error) {
	return nil, nil
}

func (dl MockNetListener) Close() error {
	return nil
}

func (dl MockNetListener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 0,
	}
}

type MockRepositoryConfigRetriever struct {
	mock.Mock
}

func (m *MockRepositoryConfigRetriever) GetConfigs(groupID int) ([]entities.Config, error) {
	args := m.Called(groupID)
	return args.Get(0).([]entities.Config), args.Error(1)
}

type MockListenerProvider struct {
	mock.Mock
}

func (m *MockListenerProvider) GetListeners(count int) ([]net.Listener, error) {
	var listeners []net.Listener
	for i := 0; i < count; i++ {
		listeners = append(listeners, &MockNetListener{})
	}
	return listeners, nil
}

type MockPortSetterFactory struct {
	mock.Mock
}

func (m *MockPortSetterFactory) Fetch(configType string) (abstraction.PortSetter, error) {
	args := m.Called(configType)
	return args.Get(0).(abstraction.PortSetter), args.Error(1)
}

type MockDomainRankFetcher struct {
	mock.Mock
}

func (m *MockDomainRankFetcher) GetDomainsWithRank(domainsCountLimit int) ([]structs.DomainWithRank, error) {
	args := m.Called(domainsCountLimit)
	return args.Get(0).([]structs.DomainWithRank), args.Error(1)
}

type MockConfigPinger struct {
	mock.Mock
}

func (m *MockConfigPinger) PingAllConfigs(configs []entities.Config, domainsWithRank []structs.DomainWithRank, wg *sync.WaitGroup, listeners []net.Listener, configScoresMap *sync.Map) {
	m.Called(configs, domainsWithRank, wg, listeners, configScoresMap)
}

type MockConfigScoreWriter struct {
	mock.Mock
}

func (m *MockConfigScoreWriter) WriteScoresToDb(configs []entities.Config, scores *sync.Map) {
	m.Called(configs, scores)
}

type MockPortSetter struct {
	mock.Mock
}

func (m *MockPortSetter) SetPort(listener net.Listener, path string) error {
	args := m.Called(listener, path)
	return args.Error(0)
}

var configOrganizerTests = []configOrganizerTest{
	{
		name:              "should call mocks and do not return error when there are multiple valid configs and domainWithRanks",
		domainsCountLimit: 5,
		configs:           []entities.Config{{Path: "config1"}, {Path: "config2"}},
		domainsWithRank: []structs.DomainWithRank{
			{Domain: entities.Domain{Address: "example.com"}, Rank: 1},
			{Domain: entities.Domain{Address: "example.org"}, Rank: 2},
			{Domain: entities.Domain{Address: "example.net"}, Rank: 3},
			{Domain: entities.Domain{Address: "example.edu"}, Rank: 4},
			{Domain: entities.Domain{Address: "example.co"}, Rank: 5},
		},
	},
	{
		name:              "should call mocks and do not return error when there are zero configs but multiple domainWithRanks",
		domainsCountLimit: 5,
		configs:           []entities.Config{}, // Zero configs
		domainsWithRank: []structs.DomainWithRank{
			{Domain: entities.Domain{Address: "example.com"}, Rank: 1},
			{Domain: entities.Domain{Address: "example.org"}, Rank: 2},
			{Domain: entities.Domain{Address: "example.net"}, Rank: 3},
			{Domain: entities.Domain{Address: "example.edu"}, Rank: 4},
			{Domain: entities.Domain{Address: "example.co"}, Rank: 5},
		},
	},
	{
		name:              "should call mocks and do not return error when there are zero domainWithRanks but multiple configs",
		domainsCountLimit: 5,
		configs:           []entities.Config{{Path: "config1"}, {Path: "config2"}}, // Multiple configs
		domainsWithRank:   []structs.DomainWithRank{},                              // Zero domains
	},
}

func successfulOrganizeTest(t *testing.T, testCase configOrganizerTest) {
	t.Parallel()
	// Arrange
	mockRetriever := new(MockRepositoryConfigRetriever)
	mockListenerProvider := new(MockListenerProvider)
	mockPortSetter := new(MockPortSetter)
	mockPortSetterFactory := new(MockPortSetterFactory)
	mockDomainRankFetcher := new(MockDomainRankFetcher)
	mockConfigPinger := new(MockConfigPinger)
	mockConfigScoreWriter := new(MockConfigScoreWriter)

	mockRetriever.On("GetConfigs", 1).Return(testCase.configs, nil)
	mockPortSetter.On("SetPort", mock.Anything, mock.Anything).Return(nil)
	mockPortSetterFactory.On("Fetch", mock.Anything).Return(mockPortSetter, nil)
	mockDomainRankFetcher.On("GetDomainsWithRank", testCase.domainsCountLimit).Return(testCase.domainsWithRank, nil)
	mockConfigPinger.On("PingAllConfigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockConfigScoreWriter.On("WriteScoresToDb", mock.Anything, mock.Anything).Return()

	organizer := usecases.NewConfigsOrganizer(mockRetriever, mockPortSetterFactory, ConfigForTest, mockConfigScoreWriter, mockConfigPinger, mockDomainRankFetcher, mockListenerProvider)

	// Act
	err := organizer.Organize(1, testCase.domainsCountLimit)

	// Assert
	assert.NoErrorf(t, err, fmt.Sprintf("%v", err))
}

func TestConfigsOrganizer(t *testing.T) {
	t.Parallel()
	for _, testCase := range configOrganizerTests {
		t.Run(testCase.name, func(t *testing.T) {
			successfulOrganizeTest(t, testCase)
		})
	}
}
