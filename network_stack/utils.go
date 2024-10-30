package networkstack

import (
	"fmt"
	"net"

	"github.com/custodiaJs/cjs-network-platform/ipnet"
)

// localProcess simuliert die Verarbeitung eines Pakets durch einen lokalen Prozess
func localProcess(packet ipnet.Packet) {
	fmt.Println("[Local Process] Handling packet:", packet)
}

// sendPacket simuliert das Senden eines Pakets
func sendPacket(packet ipnet.Packet) {
	fmt.Println("[Send Packet] Packet sent to destination:", packet.DestinationIP)
}

// isLocalAddress prüft, ob die IP-Adresse lokal ist
func isLocalAddress(ip net.IP) bool {
	localIPs := []net.IP{
		net.ParseIP("192.168.1.1"), // Beispiel IPv4-Adresse
		net.ParseIP("2001:db8::1"), // Beispiel IPv6-Adresse
	}
	for _, localIP := range localIPs {
		if localIP.Equal(ip) {
			return true
		}
	}
	return false
}
