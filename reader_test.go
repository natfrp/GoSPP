package gospp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleRead() {
	payload := []byte{86, 236, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 114, 5, 1, 4, 32, 1, 13, 184, 25, 25, 8, 16, 0, 0, 0, 0, 0, 0, 0, 43, 48, 57, 212, 49, 233, 233, 233, 233, 233, 233, 233, 233, 233}
	r := bytes.NewReader(payload)
	hdr, _ := Read(r)
	fmt.Printf("%+v\n", hdr)
	// Output: &{ClientAddr:114.5.1.4 ProxyAddr:2001:db8:1919:810::2b ClientPort:12345 ProxyPort:54321}
}

func ExampleReadBuf() {
	payload := []byte{86, 236, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 114, 5, 1, 4, 32, 1, 13, 184, 25, 25, 8, 16, 0, 0, 0, 0, 0, 0, 0, 43, 48, 57, 212, 49, 233, 233, 233, 233, 233, 233, 233, 233, 233}
	r := bytes.NewReader(payload)
	hdr, _ := Read(r)
	fmt.Printf("%+v\n", hdr)
	// Output: &{ClientAddr:114.5.1.4 ProxyAddr:2001:db8:1919:810::2b ClientPort:12345 ProxyPort:54321}
}

func TestReader(t *testing.T) {
	payload := []byte{86, 236, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 114, 5, 1, 4, 32, 1, 13, 184, 25, 25, 8, 16, 0, 0, 0, 0, 0, 0, 0, 43, 48, 57, 212, 49}
	invalid := []byte{87, 236, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 114, 5, 1, 4, 32, 1, 13, 184, 25, 25, 8, 16, 0, 0, 0, 0, 0, 0, 0, 43, 48, 57, 212, 49}
	short := []byte{86, 236, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 114, 5, 1, 4}

	Convey("Read", t, func() {
		r := bytes.NewReader(payload)
		inr := bytes.NewReader(invalid)
		srtr := bytes.NewReader(short)

		Convey("valid reader", func() {
			hdr, err := Read(r)
			So(err, ShouldBeNil)
			So(hdr.ClientAddr.String(), ShouldEqual, "114.5.1.4")
			So(hdr.ProxyAddr.String(), ShouldEqual, "2001:db8:1919:810::2b")
			So(hdr.ClientPort, ShouldEqual, 12345)
			So(hdr.ProxyPort, ShouldEqual, 54321)
		})

		Convey("valid bufio reader", func() {
			hdr, err := ReadBuf(bufio.NewReader(r))
			So(err, ShouldBeNil)
			So(hdr.ClientAddr.String(), ShouldEqual, "114.5.1.4")
			So(hdr.ProxyAddr.String(), ShouldEqual, "2001:db8:1919:810::2b")
			So(hdr.ClientPort, ShouldEqual, 12345)
			So(hdr.ProxyPort, ShouldEqual, 54321)
		})

		Convey("valid bytes.Buffer reader", func() {
			hdr, err := ReadBuf(r)
			So(err, ShouldBeNil)
			So(hdr.ClientAddr.String(), ShouldEqual, "114.5.1.4")
			So(hdr.ProxyAddr.String(), ShouldEqual, "2001:db8:1919:810::2b")
			So(hdr.ClientPort, ShouldEqual, 12345)
			So(hdr.ProxyPort, ShouldEqual, 54321)
		})

		Convey("invalid reader", func() {
			_, err := Read(inr)
			So(err, ShouldEqual, ErrNoMagic)
			_, err = Read(srtr)
			So(err, ShouldEqual, io.ErrUnexpectedEOF)
		})

		Convey("invalid bufio reader", func() {
			_, err := ReadBuf(bufio.NewReader(inr))
			So(err, ShouldEqual, ErrNoMagic)
			_, err = ReadBuf(bufio.NewReader(srtr))
			So(err, ShouldEqual, io.ErrUnexpectedEOF)
		})

		Convey("invalid bytes.Buffer reader", func() {
			_, err := ReadBuf(inr)
			So(err, ShouldEqual, ErrNoMagic)
			err = inr.UnreadByte()
			So(err.Error(), ShouldEqual, "bytes.Reader.UnreadByte: at beginning of slice")
			_, err = ReadBuf(srtr)
			So(err, ShouldEqual, io.ErrUnexpectedEOF)
			err = srtr.UnreadByte()
			So(err, ShouldBeNil)
		})
	})
}
