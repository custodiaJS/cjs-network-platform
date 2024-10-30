package usstack

import (
	"fmt"

	cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"
)

// Process führt die Postrouting-Verarbeitung durch
func (p *Postrouting) Process(packet cjsnetworkplatform.Packet, stack *UserspaceNetworkStack) {
	fmt.Println("[Postrouting] Sending packet out:", packet)
	sendPacket(packet)
}
