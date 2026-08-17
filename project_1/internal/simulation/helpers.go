package simulation

import "net/netip"

func NetworkOf(prefix netip.Prefix) netip.Prefix {
	return prefix.Masked()
}

/*
* Two interface addresses belong to the same subnet when:
* - both are valid IPv4 prefixes,
* - they use the same prefix length,
* - their masked network prefixes are equal.
*
* E.g.
* a := netip.MustParsePrefix("10.0.1.10/24")
* b := netip.MustParsePrefix("10.0.1.200/24")
* SameSubnet(a, b) // true
*
* a := netip.MustParsePrefix("10.0.1.10/24")
* b := netip.MustParsePrefix("10.0.2.10/24")
* SameSubnet(a, b) // false
 */
func SameSubnet(a, b netip.Prefix) bool {
	if !a.IsValid() || !b.IsValid() {
		return false
	}

	if !a.Addr().Is4() || !b.Addr().Is4() {
		return false
	}

	/*
	* For now, requiring identical prefix lengths keeps our model
	* deterministic and easy to reason about. We'll get into more
	* sophisticated route-prefix relationships too.
	 */
	if a.Bits() != b.Bits() {
		return false
	}

	return a.Masked() == b.Masked()
}
