package ip

import (
	"encoding/binary"

	"golang.org/x/net/ipv4"
)

// Funktion zur Berechnung der IPv4-Prüfsumme
func CalculateIPv4Checksum(header *ipv4.Header) uint16 {
	headerBytes, _ := header.Marshal()
	var sum uint32
	for i := 0; i < len(headerBytes); i += 2 {
		sum += uint32(binary.BigEndian.Uint16(headerBytes[i:]))
	}
	for (sum >> 16) > 0 {
		sum = (sum >> 16) + (sum & 0xFFFF)
	}
	return ^uint16(sum)
}
