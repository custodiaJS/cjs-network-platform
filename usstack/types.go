package usstack

import (
	"net"

	cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"
)

// Input verarbeitet lokale Pakete
type Input struct {
	next *Forward
}

// Output bereitet Pakete vor, die das System verlassen
type Output struct {
	next *Postrouting
}

// Forward leitet Pakete weiter, die nicht für das lokale System bestimmt sind
type Forward struct {
	next *Postrouting
}

// Prerouting verarbeitet eingehende Pakete
type Prerouting struct {
	next *Input
}

// Postrouting führt die endgültige Verarbeitung vor dem Senden des Pakets durch
type Postrouting struct{}

// usstack enthält die verschiedenen Verarbeitungsobjekte
type UserspaceNetworkStack struct {
	Prerouting  *Prerouting
	Input       *Input
	Forward     *Forward
	Output      *Output
	Postrouting *Postrouting
	Router      *Router
}

type RouteEntry struct {
	Subnet *net.IPNet
	Action string                             // "local" oder "forward"
	Device cjsnetworkplatform.DeviceInterface // Gerät für die Weiterleitung
}

type Router struct {
	routingTableV4 []RouteEntry
	routingTableV6 []RouteEntry
}
