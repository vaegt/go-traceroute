package traceroute

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func TestExec(test *testing.T) {
	expected := returnData()

	got := Exec(expected.Dest, expected.Timeout, expected.Tries, expected.MaxTTL, expected.Protocol, expected.Port)
	if !reflect.DeepEqual(expected, got) {
		test.Errorf("Error: Exec data expected %v got %v", expected, got)
	}
}

func TestNext(test *testing.T) {
	data := returnData()

	err := data.Next()

	if err != nil {
		if err.Error() == "dial ip4:icmp 8.8.8.8: socket: operation not permitted" {
			test.Errorf("Error: Please run this test as root\n%v", err.Error())
		} else {
			test.Errorf("Error: %v", err.Error())
		}
	}
}

func TestAll(test *testing.T) {
	data := returnData()

	err := data.All()

	if err != nil {
		if err.Error() == "dial ip4:icmp 8.8.8.8: socket: operation not permitted" {
			test.Errorf("Error: Please run this test as root\n%v", err.Error())
		} else {
			test.Errorf("Error: %v", err.Error())
		}
	}
}

func TestAllUDP(test *testing.T) {
	data := udpReturnData()

	err := data.All()

	if err != nil {
		if err.Error() == "dial ip4:icmp 8.8.8.8: socket: operation not permitted" {
			test.Errorf("Error: Please run this test as root\n%v", err.Error())
		} else {
			test.Errorf("Error: %v", err.Error())
		}
	}
}

func returnData() TraceData {
	dest := net.ParseIP("8.8.8.8")
	timeout := 1 * time.Second
	tries := 1
	maxTTL := 16
	proto := "ip4:icmp"

	return TraceData{make([][]Hop, tries), dest, timeout, tries, maxTTL, proto, 0}

}

func udpReturnData() (data TraceData) {
	data = returnData()
	data.Protocol = "udp"
	data.Port = 33434
	return
}
