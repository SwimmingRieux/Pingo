package services_test

import (
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/mock"
	"net"
	"os"
	"pingo/internal/app/services"
	"strconv"
	"testing"
	"time"
)

type MockRepositoryDomainAdder struct {
	mock.Mock
}

func (mock *MockRepositoryDomainAdder) AddDomain(address string) {
	mock.Called(address)
}

type MockPacketSource struct {
	packets chan gopacket.Packet
}

func (mock *MockPacketSource) Packets() chan gopacket.Packet {
	return mock.packets
}

func generateMockPacket(dstIP string) gopacket.Packet {
	envDstPort, _ := strconv.Atoi(os.Getenv("PINGO_DEFAULT_PORT"))
	dstPort := uint16(envDstPort)

	ethernetLayer := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x0c, 0x29, 0x3e, 0x5b, 0xf2},
		DstMAC:       net.HardwareAddr{0x00, 0x0c, 0x29, 0x3e, 0x5b, 0xf3},
		EthernetType: layers.EthernetTypeIPv4,
	}

	ipLayer := &layers.IPv4{
		SrcIP:    net.IP{127, 0, 0, 1},
		DstIP:    net.ParseIP(dstIP),
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}

	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(12345),
		DstPort: layers.TCPPort(dstPort),
	}

	err := tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	if err != nil {
		return nil
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}

	err := gopacket.SerializeLayers(buffer, options,
		ethernetLayer, ipLayer, tcpLayer,
	)
	if err != nil {
		return nil
	}

	return gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
}

type networkLogRecorderTest struct {
	name            string
	networkRequests []string
}

var testCases = []networkLogRecorderTest{
	{
		name: "should record logs when multiple requests are sent",
		networkRequests: []string{
			"192.168.1.1",
			"192.168.1.2",
			"192.168.1.3",
		},
	},
}

func addOtherTestCases() {
	bigEnough := ConfigForTest.DomainsBigEnough
	largeNetworkRequests := make([]string, 0, bigEnough*2)
	for i := 0; i < bigEnough*2; i++ {
		largeNetworkRequests = append(largeNetworkRequests, fmt.Sprintf("172.16.0.%v", i%255))
	}

	testCases = append(testCases, networkLogRecorderTest{
		name:            "should record logs when requests exceed twice the BigEnough limit",
		networkRequests: largeNetworkRequests,
	})

}

func TestRecord(t *testing.T) {
	addOtherTestCases()
	t.Parallel()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			mockRepo := new(MockRepositoryDomainAdder)
			mockPacketSource := &MockPacketSource{packets: make(chan gopacket.Packet, 10)}
			recorder := services.NewNetworkLogRecorder(mockRepo, ConfigForTest, mockPacketSource)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			expectedAddresses := testCase.networkRequests
			for _, addr := range expectedAddresses {
				mockRepo.On("AddDomain", addr).Once()
			}

			// Act
			go recorder.Record(ctx)

			for _, request := range testCase.networkRequests {
				generatedPacket := generateMockPacket(request)
				mockPacketSource.packets <- generatedPacket
			}

			time.Sleep(5 * time.Second)
			cancel()
			// Assert
			mockRepo.AssertExpectations(t)
		})
	}
}
