package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/custodiaJs/cjs-network-platform/ip"
	"github.com/custodiaJs/cjs-network-platform/tcp"
	"github.com/songgao/water"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	tcpSynFlag = 0x02
	tcpAckFlag = 0x10
	tcpFinFlag = 0x01
)

// Connection ist eine Struktur zur Verwaltung von Verbindungen
type Connection struct {
	srcIP   string
	srcPort uint16
	dstIP   string
	dstPort uint16
	iface   *water.Interface
	dataCh  chan []byte
}

// Globale Variable für den Verbindungsstatus
var connectionClosed bool = false

var connections = make(map[uint16]*Connection) // Map zur Verwaltung der Verbindungen
var connMutex sync.Mutex = sync.Mutex{}

func main() {
	// TUN-Interface konfigurieren
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = "tun0" // Der Standardname für ein TUN-Interface unter Linux

	// TUN-Interface erstellen
	iface, err := water.New(config)
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des TUN-Interfaces: %v", err)
	}
	fmt.Printf("TUN-Interface %s erstellt\n", iface.Name())

	// Konfiguriere das Interface mit ip (anstelle von ifconfig)
	if err := exec.Command("ip", "addr", "add", "10.0.0.1/24", "dev", iface.Name()).Run(); err != nil {
		log.Fatalf("Fehler beim Setzen der IP-Adresse: %v", err)
	}
	if err := exec.Command("ip", "link", "set", "dev", iface.Name(), "up").Run(); err != nil {
		log.Fatalf("Fehler beim Aktivieren des TUN-Interfaces: %v", err)
	}
	fmt.Println("TUN-Interface konfiguriert und aktiviert.")

	// Pakete lesen und verarbeiten
	readPackets(iface)
}

func readPackets(iface *water.Interface) {
	packet := make([]byte, 1500) // Puffer für eingehende Pakete
	for {
		n, err := iface.Read(packet)
		if err != nil {
			fmt.Printf("Fehler beim Lesen von Paketen: %v\n", err)
			continue
		}

		// Versuche, den IPv4-Header zu parsen
		ipv4Header, err := ipv4.ParseHeader(packet[:n])
		if err != nil {
			fmt.Printf("Ungültiges Paketformat, erster Bytewert: %d\n", packet[0])
			continue
		}

		// IPv4-Paket verarbeiten
		if err := handleIPv4Packet(ipv4Header, packet[ipv4Header.Len:n], iface); err != nil {
			fmt.Printf("Fehler beim Verarbeiten des IPv4-Pakets: %v\n", err)
		}

		// Wenn FIN-Flag empfangen, beende die Schleife
		if ipv4Header.Protocol == 6 && (packet[13]&tcpFinFlag != 0) {
			fmt.Println("TCP-Verbindung beendet, FIN-Paket empfangen.")
			break
		}
	}
}

func handleIPv4Packet(ipv4Header *ipv4.Header, payload []byte, iface *water.Interface) error {
	// Prüfen, ob die Ziel-IP 10.0.0.2 ist
	if ipv4Header.Dst.String() != "10.0.0.2" {
		return nil
	}

	// Überprüfe den Protokolltyp (ICMP oder anderes)
	switch ipv4Header.Protocol {
	case 1: // ICMP
		return handleIPv4ICMPPacket(ipv4Header, payload, iface)
	case 6:
		return handleIPv4TCPPacket(ipv4Header, payload, iface)
	default:
		fmt.Println("Unbekanntes Protokoll:", ipv4Header.Protocol)
	}
	return nil
}

func handleIPv4ICMPPacket(ipv4Header *ipv4.Header, payload []byte, iface *water.Interface) error {
	fmt.Printf("Empfangenes ICMP-Paket von %s an %s\n", ipv4Header.Src, ipv4Header.Dst)

	// ICMP-Paket parsen
	icmpPacket, err := icmp.ParseMessage(1, payload) // 1 für IPv4
	if err != nil {
		return fmt.Errorf("Fehler beim Parsen des ICMP-Pakets: %v", err)
	}

	// Überprüfen, ob es sich um einen Echo-Request handelt
	if icmpPacket.Type == ipv4.ICMPTypeEcho {
		echoRequest, ok := icmpPacket.Body.(*icmp.Echo)
		if !ok {
			return fmt.Errorf("Fehler: ICMP-Paket ist kein Echo-Request")
		}

		// Debug-Ausgabe zur Überprüfung der ICMP-Antwortdetails
		fmt.Printf("Erstelle ICMP Echo Reply von %s an %s, ID=%d, Seq=%d\n",
			ipv4Header.Dst, ipv4Header.Src, echoRequest.ID, echoRequest.Seq)

		// Erstelle ein Echo-Reply basierend auf dem Echo-Request
		icmpReply := &icmp.Message{
			Type: ipv4.ICMPTypeEchoReply, // Echo Reply (Typ 0)
			Code: 0,
			Body: &icmp.Echo{
				ID:   echoRequest.ID,
				Seq:  echoRequest.Seq,
				Data: echoRequest.Data,
			},
		}

		// ICMP Echo Reply serialisieren
		replyData, err := icmpReply.Marshal(nil)
		if err != nil {
			return fmt.Errorf("Fehler beim Erstellen des ICMP-Reply: %v", err)
		}

		// Erstelle IPv4-Header für die Antwort (mit vertauschter Quell- und Ziel-IP)
		ipv4ReplyHeader := &ipv4.Header{
			Version:  4,
			Len:      20,
			TOS:      ipv4Header.TOS,
			TotalLen: 20 + len(replyData), // IPv4-Header (20 Bytes) + ICMP-Daten
			ID:       ipv4Header.ID,
			FragOff:  0, // Keine Fragmentierung
			TTL:      64,
			Protocol: 1,              // ICMP
			Src:      ipv4Header.Dst, // Quell- und Ziel-IP vertauschen
			Dst:      ipv4Header.Src,
		}

		// IPv4-Header-Prüfsumme berechnen und setzen
		ipv4ReplyHeader.Checksum = int(ip.CalculateIPv4Checksum(ipv4ReplyHeader))

		// IPv4-Header serialisieren
		ipHeaderBytes, err := ipv4ReplyHeader.Marshal()
		if err != nil {
			return fmt.Errorf("Fehler beim Serialisieren des IPv4-Headers: %v", err)
		}

		// IPv4-Header und ICMP-Daten kombinieren
		replyPacket := append(ipHeaderBytes, replyData...)

		// Antwort über den TUN-Adapter senden
		if _, err := iface.Write(replyPacket); err != nil {
			fmt.Println("Fehler beim Senden der Antwort:", err)
		} else {
			fmt.Printf("Gesendetes ICMP-Antwortpaket von %s an %s\n", ipv4ReplyHeader.Src, ipv4ReplyHeader.Dst)
		}
	}
	return nil
}

func handleIPv4TCPPacket(ipv4Header *ipv4.Header, payload []byte, iface *water.Interface) error {
	if len(payload) < 20 {
		return fmt.Errorf("Ungültiges TCP-Paket, zu kurz")
	}

	srcPort := binary.BigEndian.Uint16(payload[0:2])
	dstPort := binary.BigEndian.Uint16(payload[2:4])
	seqNum := binary.BigEndian.Uint32(payload[4:8])
	ackNum := binary.BigEndian.Uint32(payload[8:12])
	offset := (payload[12] >> 4) * 4 // Länge des TCP-Headers in Bytes
	flags := payload[13]

	fmt.Printf("Empfangenes TCP-Paket: SrcIP=%s, SrcPort=%d, DstIP=%s, DstPort=%d, Flags=%#x\n",
		ipv4Header.Src, srcPort, ipv4Header.Dst, dstPort, flags)

	// Prüfen auf FIN-Flag
	if flags&tcpFinFlag != 0 {
		fmt.Println("TCP-Verbindung wird geschlossen, FIN-Paket empfangen.")
		connectionClosed = true
		return sendTCPResponse(ipv4Header, srcPort, dstPort, ackNum, seqNum+1, tcpAckFlag, iface)
	}

	// Prüfen auf SYN-Flag
	if flags&tcpSynFlag != 0 {
		return sendTCPResponse(ipv4Header, srcPort, dstPort, ackNum, seqNum+1, tcpSynFlag|tcpAckFlag, iface)
	}

	// Prüfen auf ACK-Flag
	if flags&tcpAckFlag != 0 && flags&tcpSynFlag == 0 {
		fmt.Println("TCP-Verbindung bestätigt. Client hat den Handshake abgeschlossen.")

		// Überprüfen, ob eine Verbindung bereits existiert
		if _, exists := connections[srcPort]; !exists {
			handleNewConnection(srcPort, dstPort, iface, ipv4Header) // Verbindung verarbeiten
		}

		// Verarbeite Anwendungsdaten
		if len(payload) > int(offset) {
			appData := payload[offset:]
			fmt.Printf("Empfangene Anwendungsdaten (%d Bytes): %s\n", len(appData), string(appData))

			// Sende die Anwendungsdaten an die Verbindung, die in handleNewConnection erstellt wurde
			conn, exists := connections[srcPort]
			if exists {
				conn.dataCh <- appData // Daten an die bestehende Verbindung senden
				fmt.Println("Daten an die Verbindung gesendet.")
			} else {
				fmt.Println("Verbindung nicht gefunden.")
			}

			// Sende ein ACK zurück
			return sendTCPResponse(ipv4Header, srcPort, dstPort, ackNum, seqNum+uint32(len(appData)), tcpAckFlag, iface)
		}
	}

	return nil
}

func sendTCPResponse(ipv4Header *ipv4.Header, srcPort, dstPort uint16, seqNum, ackNum uint32, flags byte, iface *water.Interface) error {
	// TCP-Header erstellen
	tcpHeader := make([]byte, 20)
	binary.BigEndian.PutUint16(tcpHeader[0:2], dstPort) // Quellport für die Antwort
	binary.BigEndian.PutUint16(tcpHeader[2:4], srcPort) // Zielport
	binary.BigEndian.PutUint32(tcpHeader[4:8], seqNum)  // Sequenznummer
	binary.BigEndian.PutUint32(tcpHeader[8:12], ackNum) // Acknowledgment Nummer
	tcpHeader[12] = 5 << 4                              // Headerlänge (5 32-Bit-Wörter)
	tcpHeader[13] = flags                               // Flags setzen (z.B. SYN-ACK)
	binary.BigEndian.PutUint16(tcpHeader[14:16], 65535) // Fenstergröße (maximaler Wert)

	// Berechne die TCP-Prüfsumme und setze sie im TCP-Header
	checksum := tcp.CalculateTCPChecksum(tcpHeader, ipv4Header.Dst, ipv4Header.Src)
	binary.BigEndian.PutUint16(tcpHeader[16:18], checksum)

	// IPv4-Header für die Antwort erstellen (mit vertauschter Quell- und Ziel-IP)
	ipv4ReplyHeader := &ipv4.Header{
		Version:  4,
		Len:      20,
		TOS:      ipv4Header.TOS,
		TotalLen: 20 + len(tcpHeader), // IPv4-Header (20 Bytes) + TCP-Header (20 Bytes)
		ID:       ipv4Header.ID,
		FragOff:  0,
		TTL:      64,
		Protocol: 6,              // TCP
		Src:      ipv4Header.Dst, // Quell-IP vertauschen
		Dst:      ipv4Header.Src, // Ziel-IP vertauschen
	}

	// IPv4-Header-Prüfsumme berechnen und setzen
	ipv4ReplyHeader.Checksum = int(ip.CalculateIPv4Checksum(ipv4ReplyHeader))

	// IPv4-Header serialisieren
	ipHeaderBytes, err := ipv4ReplyHeader.Marshal()
	if err != nil {
		return fmt.Errorf("Fehler beim Serialisieren des IPv4-Headers: %v", err)
	}

	// IPv4-Header und TCP-Header kombinieren und senden
	replyPacket := append(ipHeaderBytes, tcpHeader...)
	if _, err := iface.Write(replyPacket); err != nil {
		fmt.Println("Fehler beim Senden des TCP-Antwortpakets:", err)
	} else {
		fmt.Printf("Gesendetes TCP-Paket von %s an %s, SeqNum=%d, AckNum=%d, Flags=%#x\n", ipv4ReplyHeader.Src, ipv4ReplyHeader.Dst, seqNum, ackNum, flags)
	}

	return nil
}

func handleNewConnection(srcPort uint16, dstPort uint16, iface *water.Interface, ipv4Header *ipv4.Header) {
	// Erstelle ein neues Connection-Objekt
	connection := &Connection{
		srcIP:   ipv4Header.Src.String(),
		srcPort: srcPort,
		dstIP:   ipv4Header.Dst.String(),
		dstPort: dstPort,
		iface:   iface,
		dataCh:  make(chan []byte),
	}

	connections[srcPort] = connection

	// Starte eine neue Goroutine für die Verbindung
	go func(conn *Connection) {
		fmt.Println("RUN")
		for {
			t := <-conn.dataCh
			fmt.Println(t)
		}

	}(connection)
}
