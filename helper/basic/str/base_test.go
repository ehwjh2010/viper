package str

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	Convey("Str basic method", t, func() {
		var a string
		So(IsEmpty(a), ShouldBeTrue)
		So(Size(a), ShouldEqual, 0)

		a = "art"
		So(IsNotEmpty(a), ShouldBeTrue)

		a = "张三and李四|王五"
		So(Size(a), ShouldEqual, 10)
	})
}

func TestIsNotEmptySlice(t *testing.T) {
	Convey("Str slice is empty", t, func() {
		var a []string
		So(IsEmptySlice(a), ShouldBeTrue)
		a = append(a, "a")
		So(IsNotEmptySlice(a), ShouldBeTrue)
	})
}