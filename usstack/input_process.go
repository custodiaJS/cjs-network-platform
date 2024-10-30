package usstack

import (
	"fmt"

	cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"
)

// Process führt die Input-Verarbeitung durch und leitet das Paket an den lokalen Prozess weiter
func (i *Input) Process(packet cjsnetworkplatform.Packet, stack *UserspaceNetworkStack) {
	fmt.Println("[Input] Packet for local processing:", packet)
	localProcess(packet)
}
