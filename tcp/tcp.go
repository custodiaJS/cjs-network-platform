package tcp

import (
	"encoding/binary"
	"net"
)

func CalculateTCPChecksum(tcpHeader []byte, srcIP, dstIP net.IP) uint16 {
	pseudoHeader := make([]byte, 12)
	copy(pseudoHeader[0:4], srcIP.To4())
	copy(pseudoHeader[4:8], dstIP.To4())
	pseudoHeader[8] = 0
	pseudoHeader[9] = 6 // TCP-Protokollnummer
	binary.BigEndian.PutUint16(pseudoHeader[10:], uint16(len(tcpHeader)))

	sum := uint32(0)
	data := append(pseudoHeader, tcpHeader...)
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i:]))
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for (sum >> 16) > 0 {
		sum = (sum >> 16) + (sum & 0xFFFF)
	}
	return ^uint16(sum)
}
