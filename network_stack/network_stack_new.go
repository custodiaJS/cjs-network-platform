package networkstack

import "github.com/custodiaJs/cjs-network-platform/router"

// NewNetworkStack erstellt einen neuen Netzwerkstack mit allen Ketten
func NewNetworkStack() *NetworkStack {
	router := router.NewRouter()
	stack := &NetworkStack{
		prerouting:  &Prerouting{},
		input:       &Input{},
		forward:     &Forward{},
		output:      &Output{},
		postrouting: &Postrouting{},
		router:      router,
	}
	// Verbindungen zwischen den Ketten einrichten
	stack.prerouting.next = stack.input
	stack.input.next = stack.forward
	stack.forward.next = stack.postrouting
	stack.output.next = stack.postrouting
	return stack
}
