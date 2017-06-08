go-traceroute
=========
[![Go Report Card](https://goreportcard.com/badge/github.com/pl0th/go-traceroute)](https://goreportcard.com/report/github.com/pl0th/go-traceroute)

## Installation
Please be aware that mac os doesn't support the setcap command.
```bash
go get github.com/pl0th/go-traceroute
cd $GOPATH/src/github.com/pl0th/cmd
go build -o go-traceroute
sudo setcap 'cap_net_raw+p' ./go-traceroute
```

or

download the latest [release](https://github.com/pl0th/go-traceroute/releases)
