package gospp

import (
	"io"
	"net"
)

// Header only used as return value while decoding an SPP header
type Header struct {
	ClientAddr net.IP
	ProxyAddr  net.IP
	ClientPort uint16
	ProxyPort  uint16
}

// Read will first try read 2 bytes to check magic and then left 36 bytes and decode then into a Header.
func Read(r io.Reader) (*Header, error) {
	buf := make(rawHeader, sppMsgLength)

	// read the magic and checkit
	_, err := io.ReadFull(r, buf[:clientAddrOffset])
	if err != nil {
		return nil, err
	}
	if !buf.IsMagicMatch() {
		return nil, ErrNoMagic
	}

	// read the full spp header
	_, err = io.ReadFull(r, buf[clientAddrOffset:])
	if err != nil {
		return nil, err
	}

	return buf.ToStruct(), nil
}

type peeker interface {
	Peek(n int) ([]byte, error)
}

// ReadBuf will try to peak 2 bytes to check if matches magic number, and Read then.
func ReadBuf(r io.Reader) (*Header, error) {
	var peekFn func(n int) ([]byte, error)

	if rr, ok := r.(peeker); ok { // check if we can Peek(bufio.Reader)
		peekFn = rr.Peek
	} else if rr, ok := r.(io.ByteScanner); ok { // or try (Un)ReadByte(bytes.Buffer)
		peekFn = func(n int) (buf []byte, err error) {
			buf = make([]byte, n)
			for i := range buf {
				buf[i], err = rr.ReadByte()
				if err != nil {
					return nil, err
				}
			}
			for range buf {
				err = rr.UnreadByte()
				if err != nil {
					return nil, err
				}
			}
			return
		}
	} else { // or throw an error
		return nil, ErrBufferNotSupported
	}

	if magic, err := peekFn(2); err != nil {
		return nil, err
	} else {
		if magic[0] != sppMagic[0] || magic[1] != sppMagic[1] {
			return nil, ErrNoMagic
		} else {
			return Read(r)
		}
	}
}
