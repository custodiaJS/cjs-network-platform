package usstack

import (
	"fmt"

	cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"
)

// Process führt die Forward-Verarbeitung durch
func (f *Forward) Process(packet cjsnetworkplatform.Packet, stack *UserspaceNetworkStack) {
	fmt.Println("[Forward] Forwarding packet:", packet)
	f.next.Process(packet, stack)
}
