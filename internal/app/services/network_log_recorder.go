package services

import (
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"os"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/repository"
)

type NetworkLogRecorder struct {
	domainRepository repository.RepositoryDomainAdder
	configuration    *configs.Configuration
	packetSource     abstraction.PacketSource
}

func NewNetworkLogRecorder(
	domainRepository repository.RepositoryDomainAdder,
	configuration *configs.Configuration,
	packetSource abstraction.PacketSource,
) *NetworkLogRecorder {
	return &NetworkLogRecorder{
		configuration:    configuration,
		domainRepository: domainRepository,
		packetSource:     packetSource,
	}
}

func NewNetworkLogRecorderLive(domainRepository repository.RepositoryDomainAdder, configuration *configs.Configuration) (*NetworkLogRecorder, error) {
	deviceName := os.Getenv("PINGO_DEFAULT_RECORDING_DEVICE")
	mainPort := os.Getenv("PINGO_DEFAULT_PORT")

	handle, err := pcap.OpenLive(deviceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	err = handle.SetBPFFilter(fmt.Sprintf("host 127.0.0.1 and port %v", mainPort))
	if err != nil {
		handle.Close()
		return nil, err
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	return NewNetworkLogRecorder(domainRepository, configuration, packetSource), nil
}

func (recorder *NetworkLogRecorder) Record(context context.Context) {

	for {
		select {
		case packet := <-recorder.packetSource.Packets():
			networkLayer := packet.NetworkLayer()
			if networkLayer != nil {
				dst := networkLayer.NetworkFlow().Dst().String()
				recorder.domainRepository.AddDomain(dst)
			}
		case <-context.Done():
			return
		}
	}
}
