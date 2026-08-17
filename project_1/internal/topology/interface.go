package topology

import (
	"net/netip"
	"strings"
)

type InterfaceName string
type InterfaceID string

/*
* Invariants:
* Interface name cannot be empty.
* Interface ID must be unique.
* Interface names must be unique within a node. (ErrDuplicateInterfaceName - 409)
* The parent node must exist.
* The parent topology must exist.
 */

/*
 * TODOs
 * Linked interfaces require a shared subnet (Topology/Network validation)
 * For now we're only focused on valid IPv4 configurations, not network realism
 */

type Interface struct {
	ID      InterfaceID   `json:"id"`
	Name    InterfaceName `json:"name"`
	Address netip.Prefix  `json:"address"`
}

func (id InterfaceID) TrimSpace() InterfaceID {
	return InterfaceID(strings.TrimSpace(string(id)))
}

func (name InterfaceName) TrimSpace() InterfaceName {
	return InterfaceName(strings.TrimSpace(string(name)))
}
