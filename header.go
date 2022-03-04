package gospp

import "net"

// Decode will try to decode SPP Header from a bytes array, and returns new buffer with no SPP message in it.
func Decode(buf []byte) (hdr *Header, newBuf []byte, err error) {
	if len(buf) < sppMsgLength {
		err = ErrBufferShort
		return
	}
	var h rawHeader = buf[:sppMsgLength]
	if !h.IsMagicMatch() {
		err = ErrNoMagic
		return
	}

	hdr = h.ToStruct()
	if len(buf) > sppMsgLength {
		newBuf = buf[sppMsgLength:]
	}
	return
}

// rawHeader is used as a storage and decoder
type rawHeader []byte

func (h rawHeader) IsMagicMatch() bool {
	return h[0] == sppMagic[0] && h[1] == sppMagic[1]
}

func (h rawHeader) GetClientAddr() net.IP {
	ip := make(net.IP, net.IPv6len)
	copy(ip, h[clientAddrOffset:proxyAddrOffset])
	return ip
}

func (h rawHeader) GetProxyAddr() net.IP {
	ip := make(net.IP, net.IPv6len)
	copy(ip, h[proxyAddrOffset:clientPortOffset])
	return ip
}

func (h rawHeader) GetClientPort() uint16 {
	return byteOrder.Uint16(h[clientPortOffset:proxyPortOffset])
}

func (h rawHeader) GetProxyPort() uint16 {
	return byteOrder.Uint16(h[proxyPortOffset:sppMsgLength])
}

func (h rawHeader) ToStruct() *Header {
	return &Header{
		ClientAddr: h.GetClientAddr(),
		ProxyAddr:  h.GetProxyAddr(),
		ClientPort: h.GetClientPort(),
		ProxyPort:  h.GetProxyPort(),
	}
}
