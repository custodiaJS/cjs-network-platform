package router

func NewRouter() *Router {
	return &Router{
		routingTableV4: map[string]string{
			"192.168.1.0/24": "local",
			"10.0.0.0/24":    "forward",
		},
		routingTableV6: map[string]string{
			"2001:db8::/32":       "local",
			"fd00:dead:beef::/64": "forward",
		},
	}
}
