package networkstack

import (
	"fmt"

	"github.com/custodiaJs/cjs-network-platform/ipnet"
)

// Process führt die Prerouting-Verarbeitung durch und leitet das Paket weiter
func (p *Prerouting) Process(packet ipnet.Packet, stack *NetworkStack) {
	fmt.Println("[Prerouting] Packet received:", packet)
	if stack.router.ShouldForward(packet.DestinationIP) {
		stack.forward.Process(packet, stack)
	} else if isLocalAddress(packet.DestinationIP) {
		p.next.Process(packet, stack)
	} else {
		fmt.Println("No route found, dropping packet")
	}
}
