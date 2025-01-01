package services

import (
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"os"
	"pingo/internal/domain/repository"
)

type NetworkLogRecorder struct {
	domainRepository repository.DomainRepository
}

func (recorder *NetworkLogRecorder) Record(context context.Context) {
	deviceName := os.Getenv("PINGO_DEFAULT_RECORDING_DEVICE")
	mainPort := os.Getenv("PINGO_DEFAULT_PORT")
	handle, err := pcap.OpenLive(deviceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return
	}
	defer handle.Close()

	err = handle.SetBPFFilter(fmt.Sprintf("host 127.0.0.1 and port %v", mainPort))
	if err != nil {
		return
	}

	var BigEnough = 15
	addresses := make([]string, 0, BigEnough)

	defer func() {
		recorder.domainRepository.AddDomains(addresses)
	}()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case packet := <-packetSource.Packets():
			networkLayer := packet.NetworkLayer()
			if networkLayer != nil {
				dst := networkLayer.NetworkFlow().Dst().String()
				addresses = append(addresses, dst)
				if len(addresses) == BigEnough {
					recorder.domainRepository.AddDomains(addresses)
					addresses = addresses[:0]
				}
			}
		case <-context.Done():
			return
		}
	}
}
