package usstack

import (
	"fmt"

	cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"
)

// Process führt die Prerouting-Verarbeitung durch und leitet das Paket weiter
func (p *Prerouting) Process(packet cjsnetworkplatform.Packet, stack *UserspaceNetworkStack) {
	fmt.Println("[Prerouting] Packet received from device", packet.InDevice.Name())
	action, outDevice := stack.Router.GetRoute(packet.DestinationIP)
	if action == "forward" && outDevice != nil {
		packet.OutDevice = outDevice // Setze das ausgehende Gerät
		stack.Forward.Process(packet, stack)
	} else if action == "local" && isLocalAddress(packet.DestinationIP) {
		stack.Input.Process(packet, stack)
	} else {
		fmt.Println("No route found, dropping packet")
	}
}
