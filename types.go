package cjsnetworkplatform

import "net"

// IPv4Packet speichert ein vollständiges IPv4-Paket, inklusive Header und Payload
type IPv4Packet struct {
	Header  *IPv4Header
	Payload []byte
}

// IPv4Header speichert alle Felder des IPv4-Headers
type IPv4Header struct {
	VersionAndIHL  uint8
	TOS            uint8
	TotalLength    uint16
	Identification uint16
	FlagsAndOffset uint16
	TTL            uint8
	Protocol       uint8
	HeaderChecksum uint16
	SrcIP          net.IP
	DstIP          net.IP
}

// ICMPPacket speichert alle Felder eines ICMP-Pakets inklusive Daten
type ICMPPacket struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	ID       uint16
	Seq      uint16
	Data     []byte
}
