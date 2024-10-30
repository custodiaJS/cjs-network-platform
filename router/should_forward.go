package router

import (
	"fmt"
	"net"
)

// ShouldForward prüft, ob das Paket weitergeleitet werden sollte
func (r *Router) ShouldForward(ip net.IP) bool {
	if ip.To4() != nil {
		return r.checkRoutingTable(ip, r.routingTableV4)
	} else {
		return r.checkRoutingTable(ip, r.routingTableV6)
	}
}

// checkRoutingTable überprüft, ob die IP in der angegebenen Routing-Tabelle weitergeleitet werden soll
func (r *Router) checkRoutingTable(ip net.IP, table map[string]string) bool {
	for subnet, action := range table {
		_, network, err := net.ParseCIDR(subnet)
		if err != nil {
			fmt.Printf("Invalid subnet in routing table: %s\n", subnet)
			continue
		}
		if network.Contains(ip) && action == "forward" {
			return true
		}
	}
	return false
}
