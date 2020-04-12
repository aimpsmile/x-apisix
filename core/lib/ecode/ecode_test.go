package ecode

import (
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var KaixinWeiwu = New(309, "this is weiwu")

func kaixinWeiwu() error {
	return KaixinWeiwu
}
func hahaha() error {
	return errors.WithStack(kaixinWeiwu())
}
func weiwu() error {
	return errors.Wrap(hahaha(), "1.")
}
func testWrap() error {
	return errors.Wrap(weiwu(), "什么!!!!")
}

func TestSetPlt(t *testing.T) {

	Convey("TestNew", t, func() {
		err2 := weiwu()
		So(EqualError(KaixinWeiwu, err2), ShouldBeTrue)
	})

}
func TestEqualError(t *testing.T) {
	Convey("TestNew", t, func() {
		err := testWrap()
		So(EqualError(KaixinWeiwu, err), ShouldBeTrue)
	})
}
func TestNew(t *testing.T) {
	var unitList = []struct {
		code int32
		msg  string
	}{
		{
			code: 32,
			msg:  "正值",
		},
		{
			code: -32,
			msg:  "负值",
		},
		{
			code: 0,
			msg:  "零值",
		},
	}
	Convey("TestNew", t, func() {
		for _, v := range unitList {
			c := New(v.code, v.msg)
			So(c.Code(), ShouldEqual, v.code)
			So(c.Msg(), ShouldEqual, v.msg)
			So(c.Error(), ShouldEqual, v.msg)
			n := Cause(Error(c))
			So(n.Code(), ShouldEqual, c.Code())
			So(n.Msg(), ShouldEqual, c.Msg())
		}
	})
}
