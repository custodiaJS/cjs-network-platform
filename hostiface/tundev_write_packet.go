package hostiface

import cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"

func (d *TUNDevice) WritePacket(packet cjsnetworkplatform.Packet) error {
	_, err := d.ifce.Write(packet.Payload)
	return err
}
