package netpacket

func IsIPv4Packet(packet []byte) bool {
	if len(packet) == 0 {
		return false
	}
	return packet[0]>>4 == 4
}

func IsIPv6Packet(packet []byte) bool {
	if len(packet) == 0 {
		return false
	}
	return packet[0]>>4 == 6
}
