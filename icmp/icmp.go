package icmp

import (
	"encoding/binary"

	cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"
)

// Funktion zur Erstellung einer ICMP Echo Reply
func CreateICMPEchoReply(ipv4Packet *cjsnetworkplatform.IPv4Packet, icmpPacket *cjsnetworkplatform.ICMPPacket) []byte {
	// Erstelle das ICMP Echo Reply-Paket basierend auf dem empfangenen ICMP Echo Request
	icmpReply := &cjsnetworkplatform.ICMPPacket{
		Type: 0, // Echo Reply (Typ 0)
		Code: 0,
		ID:   icmpPacket.ID,
		Seq:  icmpPacket.Seq,
		Data: icmpPacket.Data,
	}
	icmpReply.Checksum = calculateICMPChecksum(icmpReply)

	// Serialisiere das ICMP Reply-Paket
	icmpData := serializeICMP(icmpReply)

	// Erstelle den IPv4-Header für die Antwort (mit vertauschter Quell- und Ziel-IP)
	ipv4ReplyHeader := &cjsnetworkplatform.IPv4Header{
		VersionAndIHL:  ipv4Packet.Header.VersionAndIHL,
		TOS:            ipv4Packet.Header.TOS,
		TotalLength:    uint16(len(icmpData) + 20), // 20 Bytes für den IPv4-Header
		Identification: ipv4Packet.Header.Identification,
		FlagsAndOffset: ipv4Packet.Header.FlagsAndOffset,
		TTL:            64,
		Protocol:       1,                       // ICMP
		SrcIP:          ipv4Packet.Header.DstIP, // Quell- und Ziel-IP vertauschen
		DstIP:          ipv4Packet.Header.SrcIP,
	}
	ipHeader := serializeIPv4Header(ipv4ReplyHeader)

	// Kombiniere IPv4-Header und ICMP-Daten zu einem vollständigen Paket
	return append(ipHeader, icmpData...)
}

// Funktion zur Berechnung der ICMP-Prüfsumme
func calculateICMPChecksum(icmp *cjsnetworkplatform.ICMPPacket) uint16 {
	data := serializeICMP(icmp)
	var sum uint32
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

// Funktion zur Serialisierung eines ICMP-Pakets
func serializeICMP(icmp *cjsnetworkplatform.ICMPPacket) []byte {
	buf := make([]byte, 8+len(icmp.Data))
	buf[0] = icmp.Type
	buf[1] = icmp.Code
	binary.BigEndian.PutUint16(buf[2:], icmp.Checksum)
	binary.BigEndian.PutUint16(buf[4:], icmp.ID)
	binary.BigEndian.PutUint16(buf[6:], icmp.Seq)
	copy(buf[8:], icmp.Data)
	return buf
}

// Funktion zur Serialisierung des IPv4-Headers
func serializeIPv4Header(header *cjsnetworkplatform.IPv4Header) []byte {
	buf := make([]byte, 20)
	buf[0] = header.VersionAndIHL
	buf[1] = header.TOS
	binary.BigEndian.PutUint16(buf[2:], header.TotalLength)
	binary.BigEndian.PutUint16(buf[4:], header.Identification)
	binary.BigEndian.PutUint16(buf[6:], header.FlagsAndOffset)
	buf[8] = header.TTL
	buf[9] = header.Protocol
	binary.BigEndian.PutUint16(buf[10:], header.HeaderChecksum)
	copy(buf[12:16], header.SrcIP.To4())
	copy(buf[16:20], header.DstIP.To4())
	return buf
}
