package cjsnetworkplatform

// Device interface repräsentiert ein Netzwerkgerät
type DeviceInterface interface {
	ReadPacket() (Packet, error)
	WritePacket(packet Packet) error
	Name() string
}
