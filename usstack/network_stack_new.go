package usstack

// NewNetworkStack erstellt einen neuen Netzwerkstack mit allen Ketten
func NewNetworkStack(router *Router) *UserspaceNetworkStack {
	stack := &UserspaceNetworkStack{
		Prerouting:  &Prerouting{},
		Input:       &Input{},
		Forward:     &Forward{},
		Output:      &Output{},
		Postrouting: &Postrouting{},
		Router:      router,
	}
	// Verbindungen zwischen den Ketten einrichten
	stack.Prerouting.next = stack.Input
	stack.Input.next = stack.Forward
	stack.Forward.next = stack.Postrouting
	stack.Output.next = stack.Postrouting
	return stack
}
