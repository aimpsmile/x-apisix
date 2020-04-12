package task

import (
	"testing"

	"github.com/micro-in-cn/x-apisix/monitor/conf"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceNameToInfo(t *testing.T) {
	var unitList = map[string]Svariable{
		"aimgo.http2.v1.passport": {
			BU:     "aimgo",
			Module: "passport",
			Stype:  "http2",
			Sver:   "v1",
		},
		"aimgo": {
			BU:     "aimgo",
			Module: "",
			Stype:  "",
			Sver:   "",
		},
		"aimgo.goods": {
			BU:     "aimgo",
			Module: "",
			Stype:  "goods",
			Sver:   "",
		},
		"aimgo.web.v1": {
			BU:     "aimgo",
			Module: "",
			Stype:  "web",
			Sver:   "v1",
		},
	}
	for k, v := range unitList {
		Convey("TestServiceNameToInfo-Case:"+k+"\n", t, func() {
			info := ServiceNameToInfo(k, "")
			So(info, ShouldNotBeNil)
			So(info.BU, ShouldEqual, v.BU)
			So(info.Module, ShouldEqual, v.Module)
			So(info.Stype, ShouldEqual, v.Stype)
			So(info.Sver, ShouldEqual, v.Sver)
			So(info.ModuleVeh, ShouldHaveSameTypeAs, (*conf.ModuleVer)(nil))
		})
	}
}

func TestShouldSnameBeConf(t *testing.T) {
	var unitList = map[string]struct {
		Svar    *Svariable
		IsMatch *conf.ModuleVer
	}{
		"aimgo.http2.v3.tian": {
			Svar: &Svariable{
				BU:      "aimgo",
				Module:  "tian",
				Stype:   "http2",
				Sver:    "v3",
				Version: "v3.2.11",
			},
			IsMatch: &conf.ModuleVer{
				RouteTpl:   "*.grpc.v3.5.*.routes.json",
				ServiceTpl: "*.grpc.v3.5.*.services.json",
				Hosts:      []string{"srv.v3.uqudu.com", "grpc.v3.uqudu.com"},
			},
		},
		"aimgo.http2.v1.passport": {
			Svar: &Svariable{
				BU:      "aimgo",
				Module:  "passport",
				Stype:   "http2",
				Sver:    "v1",
				Version: "v1.5.11",
			},
			IsMatch: &conf.ModuleVer{
				RouteTpl:   "*.grpc.*.*.routes.json",
				ServiceTpl: "*.grpc.*.*.services.json",
				Hosts:      []string{"srv.uqudu.com", "grpc.uqudu.com"},
			},
		},
		"aimgo.http2.v23.passport": {
			Svar: &Svariable{
				BU:      "aimgo",
				Module:  "passport",
				Stype:   "http2",
				Sver:    "v23",
				Version: "v2.3.11.22.3",
			},
			IsMatch: &conf.ModuleVer{
				RouteTpl:   "*.grpc.v23.*.routes.json",
				ServiceTpl: "*.grpc.v23.*.services.json",
				Hosts:      []string{"srvv23.uqudu.com"},
			},
		},
		"aimgo.web.v2.passport": {
			Svar: &Svariable{
				BU:      "aimgo",
				Module:  "passport",
				Stype:   "web",
				Sver:    "v2",
				Version: "v2.3.1",
			},
			IsMatch: &conf.ModuleVer{
				RouteTpl:   "*.web.v2.passport.routes.json",
				ServiceTpl: "*.web.v2.passport.services.json",
				Hosts:      []string{"passportv2.web.uqudu.com"},
			},
		},
		"aimgo.web.v3.passport": {
			Svar: &Svariable{
				BU:      "aimgo",
				Module:  "passport",
				Stype:   "web",
				Sver:    "v3",
				Version: "v3.0.0",
			},
			IsMatch: &conf.ModuleVer{
				RouteTpl:   "*.web.*.passport.routes.json",
				ServiceTpl: "*.web.*.passport.services.json",
				Hosts:      []string{"passport.web.uqudu.com"},
			},
		},
		"aimgo.web.v1.bbs": {
			Svar: &Svariable{
				BU:      "aimgo",
				Module:  "bbs",
				Stype:   "web",
				Sver:    "v1",
				Version: "v1.5.11",
			},
			IsMatch: &conf.ModuleVer{
				RouteTpl:   "*.web.v1.*.routes.json",
				ServiceTpl: "*.web.v1.*.services.json",
				Hosts:      []string{"webv1.uqudu.com"},
			},
		},
		"aimgo.web.v198.bbs": {
			Svar: &Svariable{
				BU:      "aimgo",
				Module:  "bbs",
				Stype:   "web",
				Sver:    "v198",
				Version: "v1.9.8",
			},
			IsMatch: &conf.ModuleVer{
				RouteTpl:   "*.web.v1.9.8.*.routes.json",
				ServiceTpl: "*.web.v1.9.8.*.services.json",
				Hosts:      []string{"webv198.uqudu.com"},
			},
		},
		"apigw": {
			Svar: &Svariable{
				BU:     "apigw",
				Module: "",
				Stype:  "",
				Sver:   "",
			},
			IsMatch: nil,
		},
		"aimgo.web": {
			Svar: &Svariable{
				BU:     "aimgo",
				Module: "",
				Stype:  "web",
				Sver:   "",
			},
			IsMatch: nil,
		},
	}
	var confList = &conf.Conf{
		Filter: []*conf.Filter{
			{
				BU:     "aimgo",
				Module: "*",
				Stype:  "http2",
				Ver: map[string]*conf.ModuleVer{
					"*": {
						RouteTpl:   "*.grpc.*.*.routes.json",
						ServiceTpl: "*.grpc.*.*.services.json",
						Hosts:      []string{"srv.uqudu.com", "grpc.uqudu.com"},
					},
					"v3.0.0,v3.6.0": {
						RouteTpl:   "*.grpc.v3.5.*.routes.json",
						ServiceTpl: "*.grpc.v3.5.*.services.json",
						Hosts:      []string{"srv.v3.uqudu.com", "grpc.v3.uqudu.com"},
					},
					"~v2.3": {
						RouteTpl:   "*.grpc.v23.*.routes.json",
						ServiceTpl: "*.grpc.v23.*.services.json",
						Hosts:      []string{"srvv23.uqudu.com"},
					},
				},
			},
			{
				BU:     "aimgo",
				Module: "*",
				Stype:  "web",
				Ver: map[string]*conf.ModuleVer{
					"*": {
						RouteTpl:   "*.web.*.*.routes.json",
						ServiceTpl: "*.web.*.*.services.json",
						Hosts:      []string{"web.uqudu.com"},
					},
					"=v1.9.8": {
						RouteTpl:   "*.web.v1.9.8.*.routes.json",
						ServiceTpl: "*.web.v1.9.8.*.services.json",
						Hosts:      []string{"webv198.uqudu.com"},
					},
					"<=v1.9.0": {
						RouteTpl:   "*.web.v1.*.routes.json",
						ServiceTpl: "*.web.v1.*.services.json",
						Hosts:      []string{"webv1.uqudu.com"},
					},
				},
			},
			{
				BU:     "aimgo",
				Module: "passport",
				Stype:  "web",
				Ver: map[string]*conf.ModuleVer{
					"*": {
						RouteTpl:   "*.web.*.passport.routes.json",
						ServiceTpl: "*.web.*.passport.services.json",
						Hosts:      []string{"passport.web.uqudu.com"},
					},
					"<v2.0.0": {
						RouteTpl:   "*.web.v1.passport.routes.json",
						ServiceTpl: "*.web.v1.passport.services.json",
						Hosts:      []string{"passportv1.web.uqudu.com"},
					},
					"v2.0.0,v2.9.9": {
						RouteTpl:   "*.web.v2.passport.routes.json",
						ServiceTpl: "*.web.v2.passport.services.json",
						Hosts:      []string{"passportv2.web.uqudu.com"},
					},
				},
			},
		},
	}
	for k, v := range unitList {
		Convey("ShouldSnameBeConf-Case:"+k+"\n", t, func() {
			info, err := ShouldSnameBeConf(k, v.Svar.Version, confList)
			if v.IsMatch != nil {
				So(err, ShouldBeNil)
				So(info, ShouldHaveSameTypeAs, (*Svariable)(nil))
				So(info.BU, ShouldEqual, v.Svar.BU)
				So(info.Module, ShouldEqual, v.Svar.Module)
				So(info.Stype, ShouldEqual, v.Svar.Stype)
				So(info.Sver, ShouldEqual, v.Svar.Sver)
				So(info.ModuleVeh, ShouldHaveSameTypeAs, (*conf.ModuleVer)(nil))
				So(info.ModuleVeh.ServiceTpl, ShouldEqual, v.IsMatch.ServiceTpl)
				So(info.ModuleVeh.RouteTpl, ShouldEqual, v.IsMatch.RouteTpl)
				So(len(info.ModuleVeh.Hosts), ShouldBeGreaterThan, 0)
				for _, h := range v.IsMatch.Hosts {
					So(h, ShouldBeIn, info.ModuleVeh.Hosts)
				}
			} else {
				t.Logf("[sname]%s[error]%v[info]%v", k, err, info)
				So(err, ShouldNotBeNil)
				So(info, ShouldBeNil)
			}
		})
	}
}
