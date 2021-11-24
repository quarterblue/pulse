package pulse

import (
	"testing"
)

func TestInitialize(t *testing.T) {
	capSize := 10
	p, nStream, err := Initialize(capSize)
	if err != nil {
		t.Fail()
	}

	if cap(nStream) != capSize {
		t.Fatalf("expected cap to be %d, but got %d", capSize, cap(nStream))
	}

	if p.Id != "id" {
		t.Fatalf("expected name to be %s, but got %s", "id", p.Id)
	}

}

func TestAddrToIdentifier(t *testing.T) {
	var iden Identifier = "172.17.49.2:9001"
	ipAddr := "172.17.49.2"
	port := "9001"
	tIden := AddrToIdentifier(ipAddr, port)

	if tIden != iden {
		t.Fatalf("expected identifier to be %s, but got %s", iden, tIden)
	}

}

func TestIdentifierToAddr(t *testing.T) {
	var iden Identifier = "172.17.49.2:9001"
	ipAddr := "172.17.49.2"
	port := "9001"

	tipAddr, tport := IdentifierToAddr(iden)
	if tipAddr != ipAddr {
		t.Fatalf("expected ipAddr to be %s, but got %s", ipAddr, tipAddr)
	}

	if tport != port {
		t.Fatalf("expected port to be %s, but got %s", port, tport)
	}

}
