package unetstack

import (
	"fmt"
	"io"
	"sync"
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

func AddKernelSpaceNetworkInterface(name string) error {
	if !InstanceIsInited() {
		return fmt.Errorf("not inited")
	}

	return nil
}

func AddUserSpaceNetworkInterface(name string) error {
	return nil
}

func InstanceIsInited() bool {
	instanceMutex.Lock()
	defer instanceMutex.Unlock()
	return instanceWasInited
}
