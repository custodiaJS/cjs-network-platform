package networkstack

import (
	"fmt"

	"github.com/custodiaJs/cjs-network-platform/ipnet"
)

// Process führt die Postrouting-Verarbeitung durch
func (p *Postrouting) Process(packet ipnet.Packet, stack *NetworkStack) {
	fmt.Println("[Postrouting] Sending packet out:", packet)
	sendPacket(packet)
}
