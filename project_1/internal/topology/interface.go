package topology

import "strings"

/*
* Invariants:
* Interface name cannot be empty.
* Interface ID must be unique.
* Interface names must be unique within a node. (ErrDuplicateInterfaceName - 409)
* The parent node must exist.
* The parent topology must exist.
 */
type Interface struct {
	ID   InterfaceID   `json:"id"`
	Name InterfaceName `json:"name"`
}

type InterfaceName string
type InterfaceID string

func (id InterfaceID) TrimSpace() InterfaceID {
	return InterfaceID(strings.TrimSpace(string(id)))
}

func (name InterfaceName) TrimSpace() InterfaceName {
	return InterfaceName(strings.TrimSpace(string(name)))
}
