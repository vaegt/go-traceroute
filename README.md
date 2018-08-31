go-traceroute
=========
[![Go Report Card](https://goreportcard.com/badge/github.com/vaegt/go-traceroute)](https://goreportcard.com/report/github.com/vaegt/go-traceroute)

## Installation
Please be aware that macOS doesn't support the setcap command.
```bash
go get github.com/vaegt/go-traceroute
cd $GOPATH/src/github.com/vaegt/go-traceroute/cmd
go build -o go-traceroute
sudo setcap 'cap_net_raw+p' ./go-traceroute
```

or

download the latest [release](https://github.com/vaegt/go-traceroute/releases)
