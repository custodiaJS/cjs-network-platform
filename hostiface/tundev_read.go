package hostiface

import cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"

func (d *TUNDevice) ReadPacket() (cjsnetworkplatform.Packet, error) {
	buf := make([]byte, 2000)
	n, err := d.ifce.Read(buf)
	if err != nil {
		return cjsnetworkplatform.Packet{}, err
	}
	data := buf[:n]
	packet := cjsnetworkplatform.Packet{
		Payload:  data,
		InDevice: d,
	}
	return packet, nil
}
