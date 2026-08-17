package tests

import (
	"net/netip"
	"netsim_1/internal/simulation"
	"testing"
)

func TestNetworkOf(t *testing.T) {
	prefix := netip.MustParsePrefix("10.0.1.37/24")

	got := simulation.NetworkOf(prefix)
	want := netip.MustParsePrefix("10.0.1.0/24")

	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestSameSubnet(t *testing.T) {
	a := netip.MustParsePrefix("10.0.1.10/24")
	b := netip.MustParsePrefix("10.0.1.200/24")

	if !simulation.SameSubnet(a, b) {
		t.Fatal("expected addresses to share a subnet")
	}
}

func TestDifferentSubnet(t *testing.T) {
	a := netip.MustParsePrefix("10.0.1.10/24")
	b := netip.MustParsePrefix("10.0.2.10/24")

	if simulation.SameSubnet(a, b) {
		t.Fatal("expected addresses to be on different subnets")
	}
}

// This one is more to enforce my own understanding of the address prefix vs. network
func TestSanityCheckPrefix(t *testing.T) {
	prefix := netip.MustParsePrefix("10.0.1.10/24")

	if prefix == prefix.Masked() {
		t.Fatal("...")
	}
}

/*
* TODO:
* valid IPv4 prefix succeeds
* invalid string is rejected
* plain IP without CIDR is rejected
* IPv6 prefix is rejected
* missing topology is rejected
* missing node is rejected
* missing interface is rejected
* existing address can be changed
 */

func TestValidIPv4Create(t *testing.T) {

}

func TestInvalidAddrCreate(t *testing.T) {

}

func TestMissingCIDRCreate(t *testing.T) {

}

func TestRejectIPv6(t *testing.T) {

}

func TestMissingTopologyRejects(t *testing.T) {

}
func TestMissingNodeRejects(t *testing.T) {

}
func TestMissingInterfaceRejects(t *testing.T) {

}

func TestUpdateExistingAddress(t *testing.T) {

}
