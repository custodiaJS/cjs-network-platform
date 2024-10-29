package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os/exec"
	"sync"

	"github.com/songgao/water"
)

func main() {
	// Erstelle das erste TUN-Interface
	ifce1, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Interface 1 Name: %s\n", ifce1.Name())

	// Erstelle das zweite TUN-Interface
	ifce2, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Interface 2 Name: %s\n", ifce2.Name())

	// Konfiguriere die Interfaces
	configureInterface(ifce1.Name(), "10.0.0.1/24")
	configureInterface(ifce2.Name(), "10.0.1.1/24")

	// Starte die Paketverarbeitung
	var wg sync.WaitGroup
	wg.Add(2)
	go handleInterface(ifce1, ifce2, &wg)
	go handleInterface(ifce2, ifce1, &wg)
	wg.Wait()
}

func configureInterface(name, ip string) {
	cmd := exec.Command("ip", "addr", "add", ip, "dev", name)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("ip", "link", "set", "dev", name, "up")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func handleInterface(ifceIn, ifceOut *water.Interface, wg *sync.WaitGroup) {
	defer wg.Done()
	packet := make([]byte, 1500)
	for {
		n, err := ifceIn.Read(packet)
		if err != nil {
			log.Fatal(err)
		}
		// Verarbeite das Paket
		handlePacket(packet[:n], ifceIn, ifceOut)
	}
}

func handlePacket(packet []byte, ifceIn, ifceOut *water.Interface) {
	// Analysiere den IP-Header
	ipHeader, payload := parseIPHeader(packet)
	if ipHeader == nil {
		return // Ungültiges Paket
	}

	// Entscheide, ob das Paket für uns ist (INPUT), weitergeleitet werden muss (FORWARD) oder von uns stammt (OUTPUT)
	if isLocalAddress(ipHeader.Destination) {
		// INPUT-Kette
		handleInput(ipHeader, payload, ifceIn)
	} else if shouldForward(ipHeader) {
		// FORWARD-Kette
		handleForward(packet, ifceOut)
	} else {
		// DROP
	}
}

func handleInput(ipHeader *IPHeader, payload []byte, ifce *water.Interface) {
	switch ipHeader.Protocol {
	case 0x01: // ICMP
		icmpType := payload[0]
		if icmpType == 8 { // Echo Request
			reply := createICMPEchoReply(ipHeader, payload)
			ifce.Write(reply)
		}
	case 0x06: // TCP
		// Handle TCP Packet
		handleTCPPacket(ipHeader, payload)
	case 0x11: // UDP
		// Handle UDP Packet
		handleUDPPacket(ipHeader, payload)
	default:
		// Unbekanntes Protokoll
	}
}

func handleForward(packet []byte, ifceOut *water.Interface) {
	// FORWARD-Kette: Paket weiterleiten
	// Verringere TTL
	packet[8] -= 1
	// Berechne die IP-Prüfsumme neu
	packet[10] = 0
	packet[11] = 0
	csum := checksum(packet[:20])
	packet[10] = byte(csum >> 8)
	packet[11] = byte(csum & 0xff)
	// Sende das Paket über das andere Interface
	ifceOut.Write(packet)
}

func isLocalAddress(ipAddr [4]byte) bool {
	// Überprüfe, ob die IP-Adresse zu unseren Interfaces gehört
	localAddresses := [][4]byte{
		{10, 0, 0, 1},
		{10, 0, 1, 1},
	}
	for _, addr := range localAddresses {
		if ipAddr == addr {
			return true
		}
	}
	return false
}

func shouldForward(ipHeader *IPHeader) bool {
	// Einfache Routing-Entscheidung: Alles weiterleiten, was nicht für uns ist
	return true
}

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

func parseIPHeader(packet []byte) (*IPHeader, []byte) {
	if len(packet) < 20 {
		return nil, nil // Ungültiges Paket
	}
	header := &IPHeader{
		Version:        packet[0] >> 4,
		IHL:            packet[0] & 0x0F,
		TotalLength:    binary.BigEndian.Uint16(packet[2:4]),
		Identification: binary.BigEndian.Uint16(packet[4:6]),
		Flags:          binary.BigEndian.Uint16(packet[6:8]),
		TTL:            packet[8],
		Protocol:       packet[9],
		Checksum:       binary.BigEndian.Uint16(packet[10:12]),
	}
	copy(header.Source[:], packet[12:16])
	copy(header.Destination[:], packet[16:20])

	ihlBytes := int(header.IHL * 4)
	if len(packet) < ihlBytes {
		return nil, nil // Ungültiges Paket
	}

	return header, packet[ihlBytes:]
}

func createICMPEchoReply(ipHeader *IPHeader, payload []byte) []byte {
	// Erstelle neuen IP-Header
	newIPHeader := make([]byte, 20)
	copy(newIPHeader, []byte{
		0x45,       // Version und IHL
		0x00,       // DSCP und ECN
		0x00, 0x00, // Total Length (wird später gesetzt)
		0x00, 0x00, // Identification
		0x00, 0x00, // Flags und Fragment Offset
		64,         // TTL
		0x01,       // Protocol (ICMP)
		0x00, 0x00, // Header Checksum (wird später berechnet)
	})
	// Tausche Quell- und Ziel-IP-Adressen
	copy(newIPHeader[12:16], ipHeader.Destination[:])
	copy(newIPHeader[16:20], ipHeader.Source[:])

	// Erstelle neuen ICMP-Header und -Payload
	icmpPacket := make([]byte, len(payload))
	copy(icmpPacket, payload)
	icmpPacket[0] = 0 // ICMP Echo Reply
	// Setze ICMP-Prüfsumme zurück
	icmpPacket[2] = 0
	icmpPacket[3] = 0
	csum := checksum(icmpPacket)
	icmpPacket[2] = byte(csum >> 8)
	icmpPacket[3] = byte(csum & 0xff)

	// Setze Total Length
	totalLength := uint16(len(newIPHeader) + len(icmpPacket))
	newIPHeader[2] = byte(totalLength >> 8)
	newIPHeader[3] = byte(totalLength & 0xff)

	// Berechne IP-Prüfsumme
	newIPHeader[10] = 0
	newIPHeader[11] = 0
	ipCsum := checksum(newIPHeader)
	newIPHeader[10] = byte(ipCsum >> 8)
	newIPHeader[11] = byte(ipCsum & 0xff)

	// Füge IP-Header und ICMP-Paket zusammen
	reply := append(newIPHeader, icmpPacket...)
	return reply
}

func checksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 + uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for (sum >> 16) > 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return ^uint16(sum)
}

// Zusätzliche Funktionen zum Handhaben von TCP- und UDP-Paketen
func handleTCPPacket(ipHeader *IPHeader, payload []byte) {
	// Parsen des TCP-Headers
	tcpHeader := parseTCPHeader(payload)
	if tcpHeader == nil {
		return
	}

	// Hier können Sie TCP-Sockets bereitstellen
	// Zum Beispiel können Sie eine Verbindung zum lokalen TCP-Server herstellen
	// und Daten senden oder empfangen

	// Hinweis: Eine vollständige Implementierung des TCP-Stacks ist sehr komplex
	// und übersteigt den Rahmen dieses Beispiels
	fmt.Println("TCP-Paket empfangen, Port:", tcpHeader.DestinationPort)
}

func handleUDPPacket(ipHeader *IPHeader, payload []byte) {
	// Parsen des UDP-Headers
	udpHeader, data := parseUDPHeader(payload)
	if udpHeader == nil {
		return
	}

	// Erstellen eines lokalen UDP-Sockets und Senden der Daten
	addr := &net.UDPAddr{
		IP:   net.IPv4(ipHeader.Destination[0], ipHeader.Destination[1], ipHeader.Destination[2], ipHeader.Destination[3]),
		Port: int(udpHeader.DestinationPort),
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Fehler beim Erstellen des UDP-Sockets:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Fehler beim Senden der UDP-Daten:", err)
		return
	}
	fmt.Println("UDP-Daten an lokalen Socket gesendet, Port:", udpHeader.DestinationPort)
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

func parseTCPHeader(data []byte) *TCPHeader {
	if len(data) < 20 {
		return nil
	}
	header := &TCPHeader{
		SourcePort:      binary.BigEndian.Uint16(data[0:2]),
		DestinationPort: binary.BigEndian.Uint16(data[2:4]),
		SequenceNumber:  binary.BigEndian.Uint32(data[4:8]),
		AckNumber:       binary.BigEndian.Uint32(data[8:12]),
		DataOffset:      data[12] >> 4,
		Flags:           binary.BigEndian.Uint16(data[12:14]) & 0x01FF,
		WindowSize:      binary.BigEndian.Uint16(data[14:16]),
		Checksum:        binary.BigEndian.Uint16(data[16:18]),
		UrgentPointer:   binary.BigEndian.Uint16(data[18:20]),
	}
	return header
}

type UDPHeader struct {
	SourcePort      uint16
	DestinationPort uint16
	Length          uint16
	Checksum        uint16
}

func parseUDPHeader(data []byte) (*UDPHeader, []byte) {
	if len(data) < 8 {
		return nil, nil
	}
	header := &UDPHeader{
		SourcePort:      binary.BigEndian.Uint16(data[0:2]),
		DestinationPort: binary.BigEndian.Uint16(data[2:4]),
		Length:          binary.BigEndian.Uint16(data[4:6]),
		Checksum:        binary.BigEndian.Uint16(data[6:8]),
	}
	return header, data[8:]
}
