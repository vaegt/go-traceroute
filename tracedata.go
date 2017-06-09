package traceroute

import (
	"net"
	"time"
)

// TraceData represents data received by executing traceroute.
type TraceData struct {
	Hops     [][]Hop
	Dest     net.IP
	Timeout  time.Duration
	Tries    int
	MaxTTL   int
	Protocol string
	Port     int
}

// Hop represents a path between a source and a destination.
type Hop struct {
	TryNumber int
	TTL       int
	AddrIP    net.IP
	AddrDNS   []string //net.IPAddr
	Latency   time.Duration
	Err       error
}
