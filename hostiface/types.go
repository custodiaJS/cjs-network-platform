package hostiface

import (
	"github.com/songgao/water"
)

// TUNDevice repräsentiert ein TUN-Interface
type TUNDevice struct {
	ifce *water.Interface
	name string
}
