package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os/exec"

	cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"
	"github.com/custodiaJs/cjs-network-platform/icmp"
	"github.com/songgao/water"
)

// TCPPacket speichert alle Felder eines TCP-Pakets inklusive Daten
type TCPPacket struct {
	Header *TCPHeader
	Data   []byte
}

// TCPHeader speichert alle Felder des TCP-Headers
type TCPHeader struct {
	SrcPort  uint16
	DstPort  uint16
	SeqNum   uint32
	AckNum   uint32
	Offset   uint8 // Data offset
	Flags    uint8
	Window   uint16
	Checksum uint16
	Urgent   uint16
}

func main() {
	// TUN-Interface konfigurieren
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = "utun8"

	// TUN-Interface erstellen
	iface, err := water.New(config)
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des TUN-Interfaces: %v", err)
	}
	fmt.Printf("TUN-Interface %s erstellt\n", iface.Name())

	// Konfiguriere das Interface mit ifconfig
	if err := exec.Command("ifconfig", iface.Name(), "10.0.0.1", "10.0.0.2", "netmask", "255.255.255.0", "up").Run(); err != nil {
		log.Fatalf("Fehler beim Hochfahren des TUN-Interfaces: %v", err)
	}
	fmt.Println("TUN-Interface konfiguriert und aktiviert.")

	// Pakete lesen und verarbeiten
	readPackets(iface)
}

// Funktion zur Verarbeitung von Paketen
func readPackets(iface *water.Interface) {
	packet := make([]byte, 1500)

	for {
		n, err := iface.Read(packet)
		if err != nil {
			fmt.Printf("Fehler beim Lesen von Paketen: %v\n", err)
			continue
		}

		// Überprüfen, ob es sich um ein IPv4-Paket handelt
		if packet[0]>>4 == 4 {
			ipv4Packet := parseIPv4Packet(packet[:n])
			// Protokolltyp überprüfen (ICMP oder TCP)
			switch ipv4Packet.Header.Protocol {
			case 1: // ICMP
				icmpPacket := parseICMPPacket(ipv4Packet.Payload)
				if icmpPacket != nil && icmpPacket.Type == 8 {
					reply := icmp.CreateICMPEchoReply(ipv4Packet, icmpPacket)
					if _, err := iface.Write(reply); err != nil {
						fmt.Println("Fehler beim Senden der Antwort:", err)
					} else {
						fmt.Println("ICMP-Echo-Antwort gesendet")
					}
				}
			case 6: // TCP
				tcpPacket := parseTCPPacket(ipv4Packet.Payload)
				if tcpPacket != nil && (tcpPacket.Header.Flags&0x02) != 0 { // SYN-Flag prüfen
					fmt.Println("TCP")
				}
			}
		}
	}
}

// Funktion zum Parsen eines IPv4-Pakets inklusive Payload
func parseIPv4Packet(packet []byte) *cjsnetworkplatform.IPv4Packet {
	headerLen := int((packet[0] & 0x0F) * 4)
	ipv4Header := &cjsnetworkplatform.IPv4Header{
		VersionAndIHL:  packet[0],
		TOS:            packet[1],
		TotalLength:    binary.BigEndian.Uint16(packet[2:4]),
		Identification: binary.BigEndian.Uint16(packet[4:6]),
		FlagsAndOffset: binary.BigEndian.Uint16(packet[6:8]),
		TTL:            packet[8],
		Protocol:       packet[9],
		HeaderChecksum: binary.BigEndian.Uint16(packet[10:12]),
		SrcIP:          net.IP(packet[12:16]),
		DstIP:          net.IP(packet[16:20]),
	}

	return &cjsnetworkplatform.IPv4Packet{
		Header:  ipv4Header,
		Payload: packet[headerLen:],
	}
}

// Funktion zum Parsen eines ICMP-Pakets
func parseICMPPacket(payload []byte) *cjsnetworkplatform.ICMPPacket {
	if len(payload) < 8 {
		return nil
	}
	return &cjsnetworkplatform.ICMPPacket{
		Type:     payload[0],
		Code:     payload[1],
		Checksum: binary.BigEndian.Uint16(payload[2:4]),
		ID:       binary.BigEndian.Uint16(payload[4:6]),
		Seq:      binary.BigEndian.Uint16(payload[6:8]),
		Data:     payload[8:],
	}
}

// Funktion zum Parsen eines TCP-Pakets
func parseTCPPacket(payload []byte) *TCPPacket {
	if len(payload) < 20 {
		return nil
	}
	tcpHeader := &TCPHeader{
		SrcPort:  binary.BigEndian.Uint16(payload[0:2]),
		DstPort:  binary.BigEndian.Uint16(payload[2:4]),
		SeqNum:   binary.BigEndian.Uint32(payload[4:8]),
		AckNum:   binary.BigEndian.Uint32(payload[8:12]),
		Offset:   (payload[12] >> 4) * 4,
		Flags:    payload[13],
		Window:   binary.BigEndian.Uint16(payload[14:16]),
		Checksum: binary.BigEndian.Uint16(payload[16:18]),
		Urgent:   binary.BigEndian.Uint16(payload[18:20]),
	}
	return &TCPPacket{
		Header: tcpHeader,
		Data:   payload[tcpHeader.Offset:],
	}
}
