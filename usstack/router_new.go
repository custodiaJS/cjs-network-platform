package usstack

import cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"

func NewRouter(devices map[string]cjsnetworkplatform.DeviceInterface) *Router {
	rtV4 := []RouteEntry{
		{
			Subnet: parseCIDR("192.168.1.0/24"),
			Action: "local",
			Device: devices["tun0"],
		},
		{
			Subnet: parseCIDR("10.0.0.0/24"),
			Action: "forward",
			Device: devices["vdev0"],
		},
	}
	rtV6 := []RouteEntry{
		{
			Subnet: parseCIDR("2001:db8::/32"),
			Action: "local",
			Device: devices["tun0"],
		},
		{
			Subnet: parseCIDR("fd00:dead:beef::/64"),
			Action: "forward",
			Device: devices["vdev0"],
		},
	}
	return &Router{
		routingTableV4: rtV4,
		routingTableV6: rtV6,
	}
}
