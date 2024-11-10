package unetstack

import (
	"sync"

	"github.com/custodiaJs/cjs-network-platform/iface/kernelspace"
	"github.com/custodiaJs/cjs-network-platform/iface/userspace"
)

var (
	instanceMutex               *sync.Mutex                            = new(sync.Mutex) // Speichert den Globalen Mutex ab
	instanceWasInited           bool                                   = false           // Gibt an ob die Instanz Initalisiert wurde
	instanceNetworkInterfaces   map[string]NIC                                           // Speichert alle Verfügbaren Netzwerk Interfaces ab
	userspaceNetworkInterfaces  map[string]*userspace.UserspaceNIC                       // Speichert alle Userspace Network Interfaces ab
	kernelspaceNetworkInerfaces map[string]*kernelspace.KernelspaceNIC                   // Speichert alle Kernel Sapce Netzwerk Interfaces ab
)
