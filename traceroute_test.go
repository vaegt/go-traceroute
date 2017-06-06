package traceroute

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func TestExec(test *testing.T) {
	expected := returnData()

	got := Exec(expected.Dest, expected.Timeout, expected.Tries, expected.MaxTTL)
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

func returnData() traceData {
	dest := net.ParseIP("8.8.8.8")
	timeout := 3 * time.Second
	tries := 4
	maxTTL := 32

	return traceData{make([][]Hop, tries), dest, timeout, tries, maxTTL}

}
