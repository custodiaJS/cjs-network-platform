package vussiface

import (
	"fmt"

	cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"
)

func NewVirtualDevice(name string) *VirtualDevice {
	return &VirtualDevice{
		name:    name,
		inChan:  make(chan cjsnetworkplatform.Packet, 100),
		outChan: make(chan cjsnetworkplatform.Packet, 100),
	}
}

func (d *VirtualDevice) ReadPacket() (cjsnetworkplatform.Packet, error) {
	packet, ok := <-d.inChan
	if !ok {
		return cjsnetworkplatform.Packet{}, fmt.Errorf("Device %s closed", d.name)
	}
	packet.InDevice = d // Setze InDevice auf dieses Gerät
	return packet, nil
}

func (d *VirtualDevice) WritePacket(packet cjsnetworkplatform.Packet) error {
	d.outChan <- packet
	return nil
}

func (d *VirtualDevice) Name() string {
	return d.name
}

// Zusätzliche Methoden zur Interaktion mit dem virtuellen Gerät
func (d *VirtualDevice) InjectPacket(packet cjsnetworkplatform.Packet) {
	d.inChan <- packet
}

func (d *VirtualDevice) ReceivePacket() (cjsnetworkplatform.Packet, error) {
	packet, ok := <-d.outChan
	if !ok {
		return cjsnetworkplatform.Packet{}, fmt.Errorf("Device %s closed", d.name)
	}
	return packet, nil
}
