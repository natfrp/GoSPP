package gospp

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleDecode() {
	payload := []byte{86, 236, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 114, 5, 1, 4, 32, 1, 13, 184, 25, 25, 8, 16, 0, 0, 0, 0, 0, 0, 0, 43, 48, 57, 212, 49, 229, 156, 168, 229, 134, 153, 230, 181, 139, 232, 175, 149}

	fmt.Println(Decode(payload))
	// Output: &{114.5.1.4 2001:db8:1919:810::2b 12345 54321} [229 156 168 229 134 153 230 181 139 232 175 149] <nil>
}

func TestDecode(t *testing.T) {
	payload := []byte{86, 236, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 114, 5, 1, 4, 32, 1, 13, 184, 25, 25, 8, 16, 0, 0, 0, 0, 0, 0, 0, 43, 48, 57, 212, 49}
	invalid := []byte{87, 236, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 114, 5, 1, 4, 32, 1, 13, 184, 25, 25, 8, 16, 0, 0, 0, 0, 0, 0, 0, 43, 48, 57, 212, 49}
	more := []byte{229, 156, 168, 229, 134, 153, 230, 181, 139, 232, 175, 149}

	Convey("For rawHeader and Decode", t, func() {
		Convey("valid payload", func() {
			hdr, b2, err := Decode(payload)
			So(hdr.ClientAddr.String(), ShouldEqual, "114.5.1.4")
			So(hdr.ProxyAddr.String(), ShouldEqual, "2001:db8:1919:810::2b")
			So(hdr.ClientPort, ShouldEqual, 12345)
			So(hdr.ProxyPort, ShouldEqual, 54321)
			So(b2, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("valid payload w/ data", func() {
			hdr, b2, err := Decode(append(payload, more...))
			So(hdr.ClientAddr.String(), ShouldEqual, "114.5.1.4")
			So(hdr.ProxyAddr.String(), ShouldEqual, "2001:db8:1919:810::2b")
			So(hdr.ClientPort, ShouldEqual, 12345)
			So(hdr.ProxyPort, ShouldEqual, 54321)
			So(b2, ShouldResemble, more)
			So(err, ShouldBeNil)
		})

		Convey("invalid magic payload", func() {
			_, _, err := Decode(invalid)
			So(err, ShouldEqual, ErrNoMagic)
		})

		Convey("short payload", func() {
			_, _, err := Decode(invalid[0:10])
			So(err, ShouldEqual, ErrBufferShort)
		})
	})
}
