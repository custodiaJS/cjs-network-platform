package networkstack

import (
	"fmt"

	"github.com/custodiaJs/cjs-network-platform/ipnet"
)

// Process führt die Input-Verarbeitung durch und leitet das Paket an den lokalen Prozess weiter
func (i *Input) Process(packet ipnet.Packet, stack *NetworkStack) {
	fmt.Println("[Input] Packet for local processing:", packet)
	localProcess(packet)
}
