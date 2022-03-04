package gospp

import (
	"encoding/binary"
	"errors"
)

var (
	sppMagic = []byte{0x56, 0xEC}

	byteOrder = binary.BigEndian
)

const (
	addrLength   = 16
	portLength   = 2
	sppMsgLength = 38

	clientAddrOffset = 2
	proxyAddrOffset  = clientAddrOffset + addrLength
	clientPortOffset = proxyAddrOffset + addrLength
	proxyPortOffset  = clientPortOffset + portLength

	v4Offset = 4
)

var (
	ErrAddrNotEligible    = errors.New("IP Addr is not eligible")
	ErrNoMagic            = errors.New("magic number not matched")
	ErrBufferShort        = errors.New("buffer to decode should be at least 38 bytes")
	ErrBufferNotSupported = errors.New("buffer to ReadBuf is not supported")
)
