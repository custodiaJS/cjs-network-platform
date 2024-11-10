package unetstack

import (
	"fmt"
	"io"

	"github.com/custodiaJs/cjs-network-platform/iface/kernelspace"
	"github.com/custodiaJs/cjs-network-platform/iface/userspace"
)

func InitUserSpaceNetworkStack(sessionSecureKey string, sessionId string) error {
	instanceMutex.Lock()
	defer instanceMutex.Unlock()

	// Es wird geprüft ob die Init Funktion bereits aufgerufen wurde
	if instanceWasInited {
		return io.EOF
	}

	// Die Maps werden erzeugt
	instanceNetworkInterfaces = make(map[string]NIC)
	userspaceNetworkInterfaces = make(map[string]*userspace.UserspaceNIC)
	kernelspaceNetworkInerfaces = make(map[string]*kernelspace.KernelspaceNIC)

	// Der Prozess wird als Initalisiert Makiert
	instanceWasInited = true

	return nil
}

func AddKernelSpaceNIC(ksiface *kernelspace.KernelspaceNIC) error {
	// Der Prozess muss für diesen Vorgang Initalisiert sein
	if !InstanceIsInited() {
		return fmt.Errorf("not inited")
	}

	// Der Mutex wird verwendet
	instanceMutex.Lock()

	// Es wird geprüft ob das Interface berits Hinzugefügt wurde+
	if _, found := instanceNetworkInterfaces[ksiface.GetID()]; found {
		instanceMutex.Unlock()
		return fmt.Errorf("NIC always registrated")
	}

	// Das Interface wird hinzugefügt
	instanceNetworkInterfaces[ksiface.GetID()] = ksiface
	kernelspaceNetworkInerfaces[ksiface.GetID()] = ksiface

	// Es wird geprüft ob es sich um

	// Der Vorgang wurde ohne Fehler beendet
	return nil
}

func AddUserSpaceNIC(usiface *userspace.UserspaceNIC) error {
	// Der Prozess muss für diesen Vorgang Initalisiert sein
	if !InstanceIsInited() {
		return fmt.Errorf("not inited")
	}

	return nil
}

func InstanceIsInited() bool {
	instanceMutex.Lock()
	defer instanceMutex.Unlock()
	return instanceWasInited
}
