package traceroute

import (
	"errors"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"time"
)

// Exec returns traceData with initialized Hops.
func Exec(dest net.IP, timeout time.Duration, tries int, maxTTL int) (data traceData) {
	return traceData{
		Hops:    make([][]Hop, tries),
		Dest:    dest,
		Timeout: timeout,
		Tries:   tries,
		MaxTTL:  maxTTL,
	}

}

// Next executes the doHop method for every try.
func (data *traceData) Next() (err error) {
	ttl := len(data.Hops[0]) + 1
	if ttl > data.MaxTTL {
		return errors.New("Maximum TTL reached")
	}
	for try := 0; try < data.Tries; try++ {
		currentHop, err := doHop(ttl, data.Dest, data.Timeout)
		if err != nil {
			return err
		}
		if currentHop.Err == nil {
			currentHop.AddrDNS, err = net.LookupAddr(currentHop.AddrIP.String()) // maybe use memoization
		}
		currentHop.TryNumber = try
		data.Hops[try] = append(data.Hops[try], currentHop)
	}
	return
}

func doHop(ttl int, dest net.IP, timeout time.Duration) (currentHop Hop, err error) {
	conn, err := net.Dial("ip4:icmp", dest.String())
	if err != nil {
		return
	}
	defer conn.Close()
	newConn := ipv4.NewConn(conn)
	if err = newConn.SetTTL(ttl); err != nil {
		return
	}
	echo := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   rand.Int(),
			Seq:  1, // TODO Sequence should be incremented every Hop & the id should be changed on every try(not random but different)
			Data: []byte("TABS"),
		}}

	req, err := echo.Marshal(nil)
	if err != nil {
		return
	}
	packetConn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return
	}
	defer packetConn.Close()

	start := time.Now()
	_, err = conn.Write(req)

	if err != nil {
		return
	}
	if err = packetConn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return
	}

	readBytes := make([]byte, 1500)                     // 1500 Bytes ethernet MTU
	_, sAddr, connErr := packetConn.ReadFrom(readBytes) // first return value (Code) might be useful

	latency := time.Since(start)

	currentHop = Hop{
		TTL:      ttl,
		Protocol: "icmp",
		Latency:  latency,
		Err:      connErr,
	}

	if connErr == nil {
		currentHop.AddrIP = net.ParseIP(sAddr.String())
		if currentHop.AddrIP == nil {
			currentHop.Err = errors.New("timeout reached")
		}
	}

	return currentHop, err
}

func (data *traceData) All() (err error) {
	for try := 0; try < data.Tries; try++ {
		for ttl := 1; ttl <= data.MaxTTL; ttl++ {
			currentHop, err := doHop(ttl, data.Dest, data.Timeout)
			if err != nil {
				return err
			}
			if currentHop.Err == nil {
				currentHop.AddrDNS, err = net.LookupAddr(currentHop.AddrIP.String()) // maybe use memoization
			}
			currentHop.TryNumber = try
			data.Hops[try] = append(data.Hops[try], currentHop)
			if currentHop.Err == nil && data.Dest.Equal(currentHop.AddrIP) {
				break
			}
		}
	}
	return
}
