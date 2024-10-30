package networkstack

import (
	"fmt"

	"github.com/custodiaJs/cjs-network-platform/ipnet"
)

// Process führt die Forward-Verarbeitung durch
func (f *Forward) Process(packet ipnet.Packet, stack *NetworkStack) {
	fmt.Println("[Forward] Forwarding packet:", packet)
	f.next.Process(packet, stack)
}
