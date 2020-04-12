package aimerror

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	var unitList = []struct {
		alist   []error
		blist   error
		isError bool
	}{
		{
			alist:   []error{errors.New("hahahah"), nil, errors.New("kaixin")},
			blist:   errors.New("hahahah"),
			isError: true,
		},
		{
			alist:   []error{nil},
			blist:   nil,
			isError: false,
		},
		{
			alist:   []error{errors.New("第一错误"), errors.New("第一错误"), nil, errors.New("第三错误")},
			blist:   errors.New("第一错误"),
			isError: true,
		},
	}
	Convey("TestAdd", t, func() {
		var errs = Errors([]error{})
		for _, e := range unitList {
			aErros := errs.Add(e.alist...)
			if e.isError {
				So(aErros.IsError(), ShouldBeTrue)
				So(aErros.GetErrors(), ShouldContain, e.blist)
			} else {
				So(aErros.IsError(), ShouldBeFalse)
			}
		}
	})

}
func TestAimErrors(t *testing.T) {
	errs := []error{errors.New("第一个错误"), errors.New("第二个错误"), errors.New("第三个错误处理")}
	Convey("TestAimErrors", t, func() {
		gErrs := Errors(errs)
		gErrs = gErrs.Add(errors.New("第四个错误"))
		gErrs = gErrs.Add(gErrs)
		So(gErrs.Error(), ShouldEqual, "第一个错误\r\n第二个错误\r\n第三个错误处理\r\n第四个错误")
	})
}

func TestIsError(t *testing.T) {
	errsTrue := []error{errors.New("第一个错误"), errors.New("第二个错误"), errors.New("第三个错误处理")}
	errsFalse := []error{}
	Convey("TestIsError", t, func() {
		t := Errors(errsTrue)
		f := Errors(errsFalse)
		So(t.IsError(), ShouldBeTrue)
		So(f.IsError(), ShouldBeFalse)
		So(t.GetErrors(), ShouldContain, errors.New("第一个错误"))
	})
}

func TestKaixin(t *testing.T) {
	addr := "aredis.service"
	if !strings.HasPrefix(addr, "redis://") {
		addr = "redis://" + addr
	}
	t.Log(addr)
}
