package traceroute

import (
	"net"
	"time"
)

type traceData struct {
	Hops    [][]Hop
	Dest    net.IP
	Timeout time.Duration
	Tries   int
	MaxTTL  int
}

// Hop represents a path between a source and a destination.
type Hop struct {
	TryNumber int
	TTL       int
	AddrIP    net.IP
	AddrDNS   []string //net.IPAddr
	Latency   time.Duration
	Protocol  string
	Err       error
}
