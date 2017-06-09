package main

import (
	"errors"
	"fmt"
	c "github.com/fatih/color"
	trace "github.com/pl0th/go-traceroute"
	"github.com/urfave/cli"
	"net"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1"
	app.Name = "go-traceroute"
	app.Usage = "A coloured traceroute implemented in golang"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "ttl, T",
			Value: 64,
			Usage: "sets the max. TTL value",
		},
		cli.Float64Flag{
			Name:  "timeout, o",
			Value: 3,
			Usage: "sets the timeout for the icmp echo request in seconds",
		},
		cli.IntFlag{
			Name:  "tries, t",
			Value: 3,
			Usage: "sets the amount of tries",
		},
		cli.StringFlag{
			Name:  "protocol, P",
			Value: "ip4:icmp",
			Usage: "sets the request protocol",
		},
		cli.IntFlag{
			Name:  "port, p",
			Value: 33434,
			Usage: "sets the port for udp requests",
		},
		cli.BoolFlag{
			Name:        "colour, c",
			Usage:       "disables colour",
			Destination: &c.NoColor,
		},
	}

	app.Action = func(ctx *cli.Context) (err error) {
		if len(ctx.Args()) == 0 {
			cli.ShowAppHelp(ctx)
			return
		}

		ip := net.ParseIP(ctx.Args()[0])

		if ip == nil {
			ips, err := net.LookupIP(ctx.Args()[0])
			if err != nil || len(ips) == 0 {
				c.Yellow("Please provide a valid IP address or fqdn")
				return cli.NewExitError(errors.New(c.RedString("Error: %v", err.Error())), 137)
			}
			ip = ips[0]
		}
		traceData := trace.TraceData{}
		if ctx.String("protocol") == "udp" {
			traceData = trace.Exec(ip, time.Duration(ctx.Float64("timeout")*float64(time.Second.Nanoseconds())), ctx.Int("tries"), ctx.Int("ttl"), ctx.String("protocol"), ctx.Int("port"))
		} else {
			traceData = trace.Exec(ip, time.Duration(ctx.Float64("timeout")*float64(time.Second.Nanoseconds())), ctx.Int("tries"), ctx.Int("ttl"), ctx.String("protocol"), 0)
		}

		hops := make([][]printData, 0)
		err = traceData.Next()
	Loop:
		for idxTry := 0; err == nil; err = traceData.Next() {
			usedIPs := make(map[string][]time.Duration)
			hops = append(hops, make([]printData, 0))
			for idx := 0; idx < traceData.Tries; idx++ {
				hop := traceData.Hops[idx][len(hops)-1]
				if len(hop.AddrDNS) == 0 {
					traceData.Hops[idx][len(hops)-1].AddrDNS = append(hop.AddrDNS, "no dns entry found")
				}

				usedIPs[hop.AddrIP.String()] = append(usedIPs[hop.AddrIP.String()], hop.Latency)
				hops[len(hops)-1] = append(hops[len(hops)-1], printData{[]time.Duration{hop.Latency}, 1, hop})
			}
			for idx := 0; idx < traceData.Tries; idx++ {
				hop := traceData.Hops[idx][len(hops)-1]
				if _, ok := usedIPs[hop.AddrIP.String()]; ok {
					addrString := fmt.Sprintf("%v (%v) ", c.YellowString(hop.AddrIP.String()), c.CyanString(hop.AddrDNS[0]))
					if hop.AddrIP == nil {
						addrString = c.RedString("no response ")
					}

					fmt.Printf("%v: %v", idxTry, addrString)
					for _, lat := range usedIPs[hop.AddrIP.String()] {
						latString, formString := lat.String(), ""
						if lat > time.Second {
							formString = fmt.Sprintf("%v ", latString[:4]+latString[len(latString)-1:])
						} else if lat < time.Millisecond && lat > time.Nanosecond {
							formString = fmt.Sprintf("%v ", latString[:4]+latString[len(latString)-3:])
						} else {
							formString = fmt.Sprintf("%v ", latString[:4]+latString[len(latString)-2:])
						}
						fmt.Printf(c.MagentaString(formString)) //µs
					}
					fmt.Println()
				}
				delete(usedIPs, hop.AddrIP.String())
				if traceData.Dest.Equal(hop.AddrIP) && traceData.Tries == idx+1 {
					break Loop
				}
			}
			idxTry++
		}
		if err != nil {
			c.Yellow("Please make sure you run this command as root")
			return cli.NewExitError(errors.New(c.RedString("Error: %v", err.Error())), 137)
		}

		return
	}

	app.Run(os.Args)

}

type printData struct {
	latencies []time.Duration
	count     int
	trace.Hop
}
