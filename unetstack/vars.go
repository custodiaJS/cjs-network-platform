package unetstack

import (
	"net"
	"sync"

	cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"
)

var (
	instanceMutex               *sync.Mutex                   = new(sync.Mutex) // Speichert den Globalen Mutex ab
	instanceWasInited           bool                          = false           // Gibt an ob die Instanz Initalisiert wurde
	instanceNetworkInterfaces   *sync.Map                                       // Speichert alle Verfügbaren Netzwerk Interfaces ab
	userspaceNetworkInterfaces  *sync.Map                                       // Speichert alle Userspace Network Interfaces ab
	kernelspaceNetworkInerfaces *sync.Map                                       // Speichert alle Kernel Sapce Netzwerk Interfaces ab
	localIPAddresses            *sync.Map                                       // Speichert alle Lokalen IP-Adressen ab
	instanzeApiSocket           *cjsnetworkplatform.ApiSocket                   // Speichert den Api Socket dieser Instanz ab
	coreSession                 net.Conn                                        // Speichert die Verbindung zum Core Prozess ab
)
