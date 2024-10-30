package cjsnetworkplatform

import "net"

// Packet definiert ein Netzwerkpaket mit grundlegenden Feldern
type Packet struct {
	SourceIP      net.IP
	DestinationIP net.IP
	Protocol      string // z. B. "TCP" oder "UDP"
	Payload       []byte
	InDevice      DeviceInterface // Gerät, von dem das Paket empfangen wurde
	OutDevice     DeviceInterface // Gerät, über das das Paket gesendet werden soll
}
