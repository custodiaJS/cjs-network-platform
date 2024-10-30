package usstack

import (
	"net"

	cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"
)

func (r *Router) GetRoute(ip net.IP) (string, cjsnetworkplatform.DeviceInterface) {
	var rt []RouteEntry
	if ip.To4() != nil {
		rt = r.routingTableV4
	} else {
		rt = r.routingTableV6
	}
	for _, entry := range rt {
		if entry.Subnet.Contains(ip) {
			return entry.Action, entry.Device
		}
	}
	return "", nil
}
