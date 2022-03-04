# GoSPP

An implementation in Go of Cloudflare's [Simple Proxy Protocol](https://developers.cloudflare.com/spectrum/reference/simple-proxy-protocol-header/).

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
    conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{
        IP:   net.IP{127, 0, 0, 1},
        Port: 2333,
        Zone: "",
    })

    // Encode an SPP message
    w := gospp.HeaderWriter{}
    w.SetClientAddr(net.ParseIP("114.5.1.4"))
    w.SetProxyAddr(net.ParseIP("191.9.8.10"))
    w.SetClientPort(12345)
    w.SetProxyPort(54321)
    w.WriteTo(conn)

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
    buf := []byte{"some spp header", "left content"}
    hdr, b2, err := gospp.Decode(buf)
    fmt.Println(hdr, b2, err)
    // b2 will be a slice of buf to left content
}
```

## LICENSE

Apache-2.0