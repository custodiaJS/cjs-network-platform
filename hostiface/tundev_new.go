package hostiface

import (
	"fmt"

	"github.com/songgao/water"
)

func NewTUNDevice(name string) (*TUNDevice, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = name
	ifce, err := water.New(config)
	if err != nil {
		return nil, err
	}
	fmt.Printf("TUN Interface Name: %s\n", ifce.Name())
	return &TUNDevice{
		ifce: ifce,
		name: ifce.Name(),
	}, nil
}
