package abstraction

import "github.com/google/gopacket"

type PacketSource interface {
	Packets() chan gopacket.Packet
}
