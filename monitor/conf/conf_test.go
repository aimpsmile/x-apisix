package conf

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConf(t *testing.T) {
	Convey("TestLoadConf-Case", t, func() {

		filepath := "/storage/code/aimgo.config/local/apisix/conf.yml"
		conf, err := loadConf(filepath)
		So(err, ShouldBeNil)
		So(len(conf.Filter), ShouldBeGreaterThan, 1)
		So(conf.Filter, ShouldHaveSameTypeAs, ([]*Filter)(nil))
		So(conf.Gateway.Baseurl, ShouldNotEqual, "")
		So(len(conf.Filter), ShouldBeGreaterThan, 0)
		So(conf.Filter[0].Stype, ShouldNotEqual, "")
		So(conf.Filter[0].BU, ShouldNotEqual, "")
		So(conf.Filter[0].Module, ShouldNotEqual, "")
		So(conf.Filter[0].Ver, ShouldNotEqual, "")
		So(len(conf.Filter[0].Ver), ShouldBeGreaterThan, 0)
		for _, f := range conf.Filter {
			for _, v := range f.Ver {
				So(len(v.Hosts), ShouldBeGreaterThan, 0)
				So(v.RouteTpl, ShouldNotEqual, "")
				So(v.ServiceTpl, ShouldNotEqual, "")
			}
		}
	})

}
func TestInitConf(t *testing.T) {
	Convey("TestInitConf-Case", t, func() {

		filepath := "/storage/code/aimgo.config/local/apisix/conf.yml"
		err := InitConf(filepath)
		conf := MConf()
		So(err, ShouldBeNil)
		So(len(conf.Filter), ShouldBeGreaterThan, 1)
		So(conf.Filter, ShouldHaveSameTypeAs, ([]*Filter)(nil))
		So(conf.Gateway.Baseurl, ShouldNotEqual, "")
		So(len(conf.Filter), ShouldBeGreaterThan, 0)
		So(conf.Filter[0].Stype, ShouldNotEqual, "")
		So(conf.Filter[0].BU, ShouldNotEqual, "")
		So(conf.Filter[0].Module, ShouldNotEqual, "")
		So(conf.Filter[0].Ver, ShouldNotEqual, "")
		So(len(conf.Filter[0].Ver), ShouldBeGreaterThan, 0)
		for _, f := range conf.Filter {
			for _, v := range f.Ver {
				So(len(v.Hosts), ShouldBeGreaterThan, 0)
				So(v.RouteTpl, ShouldNotEqual, "")
				So(v.ServiceTpl, ShouldNotEqual, "")
			}
		}
	})

}
