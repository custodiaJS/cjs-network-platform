package router

type RoutingTable map[string]string

type Router struct {
	routingTableV4 RoutingTable
	routingTableV6 RoutingTable
}
