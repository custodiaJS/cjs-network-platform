package networkstack

import "github.com/custodiaJs/cjs-network-platform/router"

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

// NetworkStack enthält die verschiedenen Verarbeitungsobjekte
type NetworkStack struct {
	prerouting  *Prerouting
	input       *Input
	forward     *Forward
	output      *Output
	postrouting *Postrouting
	router      *router.Router
}
