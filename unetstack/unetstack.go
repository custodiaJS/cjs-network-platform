package unetstack

import (
	"fmt"
	"io"
	"sync"

	"github.com/custodiaJs/cjs-network-platform/iface/kernelspace"
	"github.com/custodiaJs/cjs-network-platform/iface/userspace"
)

func InitUserSpaceNetworkStack(sessionSecureKey string, sessionId string) error {
	instanceMutex.Lock()
	defer instanceMutex.Unlock()

	if instanceWasInited {
		return io.EOF
	}

	instanceNetworkInterfaces = new(sync.Map)
	userspaceNetworkInterfaces = new(sync.Map)
	kernelspaceNetworkInerfaces = new(sync.Map)
	localIPAddresses = new(sync.Map)

	instanceWasInited = true

	return nil
}

func AddKernelSpaceNIC(name string, ksiface *kernelspace.KernelspaceNIC) error {
	if !InstanceIsInited() {
		return fmt.Errorf("not inited")
	}

	return nil
}

func AddUserSpaceNIC(name string, usiface *userspace.UserspaceNIC) error {
	return nil
}

func InstanceIsInited() bool {
	instanceMutex.Lock()
	defer instanceMutex.Unlock()
	return instanceWasInited
}
