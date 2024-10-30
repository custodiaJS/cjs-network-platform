package main

import (
	"log"
	"net"
	"sync"

	cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"
	"github.com/custodiaJs/cjs-network-platform/usstack"
	"github.com/custodiaJs/cjs-network-platform/vussiface"
)

func main() {
	// Geräte erstellen
	devices := make(map[string]cjsnetworkplatform.DeviceInterface)

	// Virtuelles Gerät erstellen
	virtualDevice := vussiface.NewVirtualDevice("vdev0")
	devices["vdev0"] = virtualDevice

	// Router und Netzwerkstack erstellen
	router := usstack.NewRouter(devices)
	stack := usstack.NewNetworkStack(router)

	// WaitGroup erstellen
	var wg sync.WaitGroup

	// Anzahl der zu verarbeitenden Pakete
	packetCount := 1

	// Starten des Lesens von Geräten
	for _, dev := range devices {
		wg.Add(1)
		go func(d cjsnetworkplatform.DeviceInterface) {
			defer wg.Done()
			for i := 0; i < packetCount; i++ {
				packet, err := d.ReadPacket()
				if err != nil {
					log.Printf("Error reading from device %s: %v", d.Name(), err)
					continue
				}
				// Paket verarbeiten
				stack.Prerouting.Process(packet, stack)
			}
		}(dev)
	}

	// Zum Testen Pakete in das virtuelle Gerät injizieren
	for i := 0; i < packetCount; i++ {
		testPacket := cjsnetworkplatform.Packet{
			SourceIP:      net.ParseIP("192.168.1.100"),
			DestinationIP: net.ParseIP("10.0.0.5"),
			Protocol:      "UDP",
			Payload:       []byte("Test data"),
			InDevice:      virtualDevice,
		}
		virtualDevice.InjectPacket(testPacket)
	}

	// Warten, bis alle Goroutinen fertig sind
	wg.Wait()
}
