package traceroute

import (
	"errors"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// Exec returns TraceData with initialized Hops and inserts the IP version into the protocol
func Exec(dest net.IP, timeout time.Duration, tries int, maxTTL int, proto string, port int) (data TraceData) {
	data = TraceData{
		Hops:    make([][]Hop, tries),
		Dest:    dest,
		Timeout: timeout,
		Tries:   tries,
		MaxTTL:  maxTTL,
		Port:    port,
		Proto:   proto,
	}
	if dest.To4() == nil {
		data.IPv = "6"
	} else {
		data.IPv = "4"
	}
	return
}

// Next executes the doHop method for every try.
func (data *TraceData) Next() (err error) {
	ttl := len(data.Hops[0]) + 1
	if ttl > data.MaxTTL {
		return errors.New("Maximum TTL reached")
	}
	for try := 0; try < data.Tries; try++ {
		currentHop, err := doHop(ttl, data.Dest, data.Timeout, data.Proto, data.Port, data.IPv)
		if err != nil {
			return err
		}
		if currentHop.Err == nil {
			currentHop.AddrDNS, _ = net.LookupAddr(currentHop.AddrIP.String()) // maybe use memoization
		}
		currentHop.TryNumber = try
		data.Hops[try] = append(data.Hops[try], currentHop)
	}
	return
}

func doHop(ttl int, dest net.IP, timeout time.Duration, proto string, port int, ipv string) (currentHop Hop, err error) {
	var destString string
	if port == 0 {
		destString = dest.String()
	} else {
		destString = dest.String() + ":" + strconv.Itoa(port)
	}
	req := []byte{}
	dialProto := proto

	if proto == "udp" {
		req = []byte("TABS")
		dialProto += ipv
	} else if proto == "icmp" {
		dialProto = "ip" + ipv + ":" + proto
	} else {
		return currentHop, errors.New("protocol not implemented")
	}

	conn, err := net.Dial(dialProto, destString)
	if err != nil {
		return
	}
	defer conn.Close()

	listenAddress := "0.0.0.0"

	if ipv == "4" {
		newConn := ipv4.NewConn(conn)
		if err = newConn.SetTTL(ttl); err != nil {
			return
		}
		if proto == "icmp" {
			req, err = createICMPEcho(ipv4.ICMPTypeEcho)
			if err != nil {
				return
			}
		}
	} else if ipv == "6" {
		listenAddress = "::0"
		newConn := ipv6.NewConn(conn)
		if err = newConn.SetHopLimit(ttl); err != nil {
			return
		}
		if proto == "icmp" {
			req, err = createICMPEcho(ipv6.ICMPTypeEchoRequest)
			if err != nil {
				return
			}
		}
	}

	packetConn, err := icmp.ListenPacket("ip"+ipv+":"+"icmp", listenAddress)
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
		TTL:     ttl,
		Latency: latency,
		Err:     connErr,
	}

	if connErr == nil {
		currentHop.AddrIP = net.ParseIP(sAddr.String())
		if currentHop.AddrIP == nil {
			currentHop.Err = errors.New("timeout reached")
		}
	}

	return currentHop, err
}

// All executes all doHops for all tries.
func (data *TraceData) All() (err error) {
	for try := 0; try < data.Tries; try++ {
		for ttl := 1; ttl <= data.MaxTTL; ttl++ {
			currentHop, err := doHop(ttl, data.Dest, data.Timeout, data.Proto, data.Port, data.IPv)
			if err != nil {
				return err
			}
			if currentHop.Err == nil {
				currentHop.AddrDNS, _ = net.LookupAddr(currentHop.AddrIP.String()) // maybe use memoization
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

func createICMPEcho(ICMPTypeEcho icmp.Type) (req []byte, err error) {
	echo := icmp.Message{
		Type: ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   rand.Int(),
			Seq:  1, // TODO Sequence should be incremented every Hop & the id should be changed on every try(not random but different)
			Data: []byte("TABS"),
		}}

	req, err = echo.Marshal(nil)
	return
}
