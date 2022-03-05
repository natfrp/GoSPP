package gospp

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHeaderWriter(t *testing.T) {
	Convey("For a HeaderWriter", t, func() {
		w := HeaderWriter{}
		buf := &bytes.Buffer{}

		Convey("illegal IP", func() {
			err := w.SetClientAddr(nil)
			So(err, ShouldEqual, ErrAddrNotEligible)
		})

		Convey("magic number", func() {
			buf.Reset()
			w.WriteTo(buf)
			So(buf.Bytes()[:clientAddrOffset], ShouldResemble, sppMagic)
		})

		Convey("IP and v4 padding", func() {
			buf.Reset()
			w.SetClientAddr(net.ParseIP("114.5.1.4"))
			w.SetProxyAddr(net.ParseIP("191.9.8.10"))
			w.WriteTo(buf)
			So(buf.Bytes()[clientAddrOffset:proxyAddrOffset], ShouldResemble, []byte(net.ParseIP("114.5.1.4").To16()))
			So(buf.Bytes()[proxyAddrOffset-4:proxyAddrOffset], ShouldResemble, []byte(net.ParseIP("114.5.1.4").To4()))
			So(buf.Bytes()[proxyAddrOffset:clientPortOffset], ShouldResemble, []byte(net.ParseIP("191.9.8.10").To16()))
			So(buf.Bytes()[clientPortOffset-4:clientPortOffset], ShouldResemble, []byte(net.ParseIP("191.9.8.10").To4()))
		})

		Convey("ports", func() {
			buf.Reset()
			w.SetClientPort(12345)
			w.SetProxyPort(54321)
			w.WriteTo(buf)
			So(byteOrder.Uint16(buf.Bytes()[clientPortOffset:proxyPortOffset]), ShouldEqual, uint16(12345))
			So(byteOrder.Uint16(buf.Bytes()[proxyPortOffset:sppMsgLength]), ShouldEqual, uint16(54321))
		})
	})
}

func ExampleHeaderWriter() {
	w := HeaderWriter{}

	w.SetClientAddr(net.ParseIP("114.5.1.4"))
	w.SetClientPort(12345)
	w.SetProxyAddr(net.ParseIP("2001:db8:1919:810::2b"))
	w.SetProxyPort(54321)

	buf := &bytes.Buffer{}
	w.WriteTo(buf)
	fmt.Println(buf.Bytes())
}

func ExampleHeaderWriter_asSlice() {
	w := HeaderWriter{}

	w.SetClientAddr(net.ParseIP("114.5.1.4"))
	w.SetClientPort(12345)
	w.SetProxyAddr(net.ParseIP("2001:db8:1919:810::2b"))
	w.SetProxyPort(54321)

	buf := &bytes.Buffer{}
	buf.Write(w[:])
	fmt.Println(buf.Bytes())
}

func ExampleHeaderWriter_withPool() {
	sppPool := sync.Pool{New: func() interface{} {
		return &HeaderWriter{}
	}}

	w := sppPool.Get().(*HeaderWriter)

	w.SetClientAddr(net.ParseIP("114.5.1.4"))
	w.SetClientPort(12345)
	w.SetProxyAddr(net.ParseIP("2001:db8:1919:810::2b"))
	w.SetProxyPort(54321)

	buf := &bytes.Buffer{}
	w.WriteTo(buf)
	sppPool.Put(w)
	fmt.Println(buf.Bytes())
}
