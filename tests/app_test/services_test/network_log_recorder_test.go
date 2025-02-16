package services_test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"pingo/configs"
	"pingo/internal/app/services"
	"testing"
	"time"
)

type MockRepositoryDomainAdder struct {
	mock.Mock
}

func (mock *MockRepositoryDomainAdder) AddDomains(addresses []string) {
	for _, addr := range addresses {
		mock.Called(addr)
	}
}

type networkLogRecorderTest struct {
	name            string
	areParallel     bool
	networkRequests []string
}

var configForNetworkLogRecorderTest, _ = configs.NewConfig()

var testCases = []networkLogRecorderTest{
	{
		name:        "should record logs when multiple sequential requests are sent",
		areParallel: false,
		networkRequests: []string{
			"192.168.1.1",
			"192.168.1.2",
			"192.168.1.3",
		},
	},
	{
		name:        "should record logs when parallel requests are sent",
		areParallel: true,
		networkRequests: []string{
			"10.0.0.1",
			"10.0.0.2",
			"10.0.0.3",
		},
	},
}

func addOtherTestCases() {
	bigEnough := configForNetworkLogRecorderTest.DomainsBigEnough
	largeNetworkRequests := make([]string, 0, bigEnough*2)
	for i := 0; i < bigEnough*2; i++ {
		largeNetworkRequests = append(largeNetworkRequests, "172.16.0."+string(rune(i%255)))
	}

	testCases = append(testCases, networkLogRecorderTest{
		name:            "should record logs when requests exceed twice the BigEnough limit and sent consequentially",
		areParallel:     false,
		networkRequests: largeNetworkRequests,
	})

	testCases = append(testCases, networkLogRecorderTest{
		name:            "should record logs when requests exceed twice the BigEnough limit and sent in parallel",
		areParallel:     true,
		networkRequests: largeNetworkRequests,
	})

}

func TestRecord(t *testing.T) {
	addOtherTestCases()
	t.Parallel()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Assert
			mockRepo := new(MockRepositoryDomainAdder)
			recorder := services.NewNetworkLogRecorder(mockRepo, configForNetworkLogRecorderTest)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			mockRepo.On("AddDomains", mock.Anything).Return()

			expectedAddresses := testCase.networkRequests
			for _, addr := range expectedAddresses {
				mockRepo.On("AddDomains", addr).Once()
			}

			// Act
			go recorder.Record(ctx)

			// Simulate network requests (requires actual packet injection setup)

			// Assert
			time.Sleep(1 * time.Second)
			mockRepo.AssertExpectations(t)
		})
	}
}

// todo: send actual requests
// todo: don't use time.Sleep if you can
// todo: how to use the areParallel thing?
// todo: what happens to context? shouldn't we close it?
// todo: think about how are we using t.parallel and is it correct?
// todo: place realConfig in one file and use it everywhere
