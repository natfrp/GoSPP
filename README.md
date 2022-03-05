# GoSPP

[![Go Reference](https://pkg.go.dev/badge/github.com/natfrp/gospp.svg)](https://pkg.go.dev/github.com/natfrp/gospp)

An implementation in Go of
Cloudflare's [Simple Proxy Protocol](https://developers.cloudflare.com/spectrum/reference/simple-proxy-protocol-header/)
.

## Usage

```go
package main

import (
    "bufio"
    "fmt"
    "net"

    "github.com/natfrp/gospp"
)

func main() {
    // Encode an SPP message
    w := gospp.HeaderWriter{}
    w.SetClientAddr(net.ParseIP("114.5.1.4"))
    w.SetProxyAddr(net.ParseIP("191.9.8.10"))
    w.SetClientPort(12345)
    w.SetProxyPort(54321)
    w.WriteTo(conn)

    // Encode result as a slice
    buf := make([]byte, 1024)
    w.SetProxyAddr(listener.LocalAddr().(*net.UDPAddr).IP)
    w.SetProxyPort(uint16(listener.LocalAddr().(*net.UDPAddr).Port))
    w.SetClientAddr(someIP)
    w.SetClientPort(somePort)
    buf = append(w[:], buf...)

    // Decode an SPP message from reader
    hdr, err := gospp.Read(conn)
    if err != nil {
        return
    }
    fmt.Println(hdr)

    // Decode an SPP message from buffered reader
    hdr, err = gospp.ReadBuf(bufio.NewReader(conn))
    if err != nil {
        return
    }
    fmt.Println(hdr)

    // Decode an SPP message from buffer and get left content
    buf = []byte{"some spp header", "left content"}
    hdr, b2, err := gospp.Decode(buf)
    fmt.Println(hdr, b2, err)
    // b2 will be a slice of buf to left content
}
```

## LICENSE

Apache-2.0
