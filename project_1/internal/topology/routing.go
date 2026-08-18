package topology

import "net/netip"

/**
 * Future TODO's for realistic simulations
 * - Real routing systems can support recursive next-hop resolution, we don't.
 *   A static route's next hop must be directly reachable through its specified outgoing interface.
 *
 * - Real routers can have multiple routes to the same prefix and select between them using metrics,
 *   administrative distance, ECMP, routing protocols, etc. This is not today's problem to solve.
 */

/**
* Architectural decision: Don't store connected routes in Nodes
* If a node contains eth0.Address = 10.0.1.10/24
* Then we automatically know that it has access to 10.0.1.0/24
* Don't store both; derive 10.0.1.0/24 when its asked for.
*
* When somebody asks for a node's effective routing table, we construct:
*
* [connected routes derived from interfaces] + [configured static routes]
*
* This will become particularly useful for:
* - interface address changes
* - interface deletion
* - interface disablement
* because we won't have stale routing-table state to synchronize.
* This principle is present throughout this software:
* 	Store authoritative state; derive what can safely be derived.
 */

type RouteSource string

const (
	RouteSourceConnected RouteSource = "connected"
	RouteSourceStatic    RouteSource = "static"
)

/**
* The NextHop pointer has useful semantics here.
*
* For a directly connected route:
*
* 	Route{
* 	    Destination: netip.MustParsePrefix("10.0.1.0/24"),
* 	    NextHop:     nil,
* 	    InterfaceID: "eth0",
* 	    Source:      RouteSourceConnected,
* 	}
*
* We're saying:
*
* 10.0.1.0/24 is directly reachable through eth0.
*
* For a static route:
*
* gateway := netip.MustParseAddr("10.0.1.1")
*
* 	Route{
* 	    Destination: netip.MustParsePrefix("0.0.0.0/0"),
* 	    NextHop:     &gateway,
* 	    InterfaceID: "eth0",
* 	    Source:      RouteSourceStatic,
* 	}
*
* We're saying:
*
* Anything else -> give it to 10.0.1.1 (a router/gateway) through eth0.
 */
type Route struct {
	Destination         netip.Prefix `json:"destination"`
	NextHop             *netip.Addr  `json:"next_hop,omitempty"`
	InterfaceIdentifier InterfaceID  `json:"interface_id"`
	Source              RouteSource  `json:"source"`
}

/*
*
* Type distinctions:
*
* Node.StaticRoutes = user configuration
* Route = result used by routing logic
 */
type StaticRoute struct {
	Destination netip.Prefix
	NextHop     netip.Addr
	InterfaceID InterfaceID
}

/*
* Look for a compatible destination represented with the
* largest number of bits in the CIDR
 */
func BestRoute(routes []Route, destination netip.Addr) (Route, bool) {
	var best Route
	bestBits := -1

	for _, route := range routes {
		if !route.Destination.Contains(destination) {
			continue
		}

		bits := route.Destination.Bits()

		if bits > bestBits {
			best = route
			bestBits = bits
		}
	}

	if bestBits == -1 {
		return Route{}, false
	}

	return best, true
}
