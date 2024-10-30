package ipnet

type IPHeader struct {
	Version        uint8
	IHL            uint8
	TotalLength    uint16
	Identification uint16
	Flags          uint16
	TTL            uint8
	Protocol       uint8
	Checksum       uint16
	Source         [4]byte
	Destination    [4]byte
}

type TCPHeader struct {
	SourcePort      uint16
	DestinationPort uint16
	SequenceNumber  uint32
	AckNumber       uint32
	DataOffset      uint8
	Flags           uint16
	WindowSize      uint16
	Checksum        uint16
	UrgentPointer   uint16
}

type UDPHeader struct {
	SourcePort      uint16
	DestinationPort uint16
	Length          uint16
	Checksum        uint16
}
