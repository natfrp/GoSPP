package gospp

import (
	"io"
	"net"
)

var v4Padding = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}

// HeaderWriter can be reused using a pool to reduce memory allocation.
// SPP header is 38 byte fixed-length, so data being marshalled into a 38 byte array.
type HeaderWriter [sppMsgLength]byte

func writeAddr(dst []byte, addr net.IP) error {
	switch len(addr) {
	case net.IPv6len:
		copy(dst, addr)
	case net.IPv4len:
		copy(dst[:len(dst)-v4Offset], v4Padding)
		copy(dst[len(dst)-v4Offset:], addr)
	default:
		return ErrAddrNotEligible
	}
	return nil
}

// SetClientAddr sets address of the originator of the proxied UDP datagram, i.e. the client.
func (w *HeaderWriter) SetClientAddr(addr net.IP) error {
	return writeAddr(w[clientAddrOffset:proxyAddrOffset], addr)
}

// SetProxyAddr sets address of the recipient of the proxied UDP datagram, i.e. the proxy.
func (w *HeaderWriter) SetProxyAddr(addr net.IP) error {
	return writeAddr(w[proxyAddrOffset:clientPortOffset], addr)
}

// SetClientPort sets source port number of the proxied UDP datagram.
// In other words, the UDP port number from which the client sent the datagram.
func (w *HeaderWriter) SetClientPort(port uint16) {
	byteOrder.PutUint16(w[clientPortOffset:proxyPortOffset], port)
}

// SetProxyPort sets destination port number of the proxied UDP datagram.
// In other words, the UDP port number on which the proxy received the datagram.
func (w *HeaderWriter) SetProxyPort(port uint16) {
	byteOrder.PutUint16(w[proxyPortOffset:sppMsgLength], port)
}

// WriteTo implements io.WriterTo
func (w *HeaderWriter) WriteTo(wri io.Writer) (int64, error) {
	w[0], w[1] = 0x56, 0xEC
	n, err := wri.Write(w[:])
	return int64(n), err
}
