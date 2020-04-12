package core

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/micro-in-cn/x-apisix/core/logger/zaplog"
	"github.com/micro-in-cn/x-apisix/core/testconf"

	"github.com/micro-in-cn/x-apisix/core/config"
	. "github.com/smartystreets/goconvey/convey"
)

type TestList struct {
	aimEnv          string
	aimConfig       string
	sName           string
	regCenterPrefix string
	confPath        []string
	eq              string
	isErr           bool
}

func TestDoAimEnv(t *testing.T) {
	var unitList = []TestList{
		{
			aimEnv: config.EnvLocal,
			eq:     config.EnvLocal,
			isErr:  false,
		},
		{
			aimEnv: strings.ToUpper(config.EnvTest),
			eq:     config.EnvTest,
			isErr:  false,
		},
		{
			aimEnv: "abc",
			eq:     "",
			isErr:  true,
		},
	}
	Convey("TestDoAimEnv", t, func() {
		for _, v := range unitList {
			os.Setenv("AIM_ENV", v.aimEnv)
			aimEnv, err := doAimEnv()
			if v.isErr {
				So(err, ShouldBeError)
			} else {
				So(err, ShouldBeNil)
			}
			So(aimEnv, ShouldEqual, v.eq)

		}

	})
}

func TestDoAimConfig(t *testing.T) {
	aimEnv := config.EnvLocal
	confDir := testconf.GetTestConfPath()
	var unitList = []TestList{
		{
			aimConfig: confDir,
			eq:        fmt.Sprintf("%s/%s", confDir, aimEnv),
			isErr:     false,
		},
		{
			aimConfig: "/usr/local/",
			eq:        "",
			isErr:     true,
		},
		{
			aimConfig: "/tmp/",
			eq:        "",
			isErr:     true,
		},
	}
	Convey("TestDoAimConfig", t, func() {
		for _, v := range unitList {
			os.Setenv("AIM_CONFIG", v.aimConfig)
			//os.Setenv("MICRO_SERVER_NAME", v.sName)
			aimConfig, err := doAimConfig(aimEnv)
			if v.isErr {
				So(err, ShouldBeError)
			} else {
				So(err, ShouldBeNil)
			}
			So(aimConfig, ShouldEqual, v.eq)

		}

	})
}

func TestConfFile(t *testing.T) {

	confDir := testconf.GetTestConfPath()
	aimEnv := config.EnvLocal
	aimConfig := fmt.Sprintf("%s/%s", confDir, aimEnv)
	var unitList = []TestList{
		{
			confPath: []string{"db.yml"},
			isErr:    false,
		},
		{
			confPath: []string{"kaixin.yml"},
			isErr:    true,
		},
	}
	Convey("TestconfFile", t, func() {
		So(AimC(), ShouldNotBeNil)
		for _, v := range unitList {
			err := confFile(aimConfig, v.confPath...)
			if v.isErr {
				So(err, ShouldBeError)
				t.Log(err)
			} else {
				So(err, ShouldBeNil)
				So(AimC().Register, ShouldNotBeEmpty)
				So(AimC().Depend, ShouldNotBeEmpty)
			}
		}

	})
}

func TestRegCenter(t *testing.T) {
	type testConf struct {
		Level      string `json:"level"`
		LogFileDir string `json:"logFileDir"`
	}
	confDir := testconf.GetTestConfPath()
	aimEnv := config.EnvLocal
	aimConfig := fmt.Sprintf("%s/%s", confDir, aimEnv)
	err := confFile(aimConfig, "db.yml", "srv.yml")

	var unitList = []TestList{
		{
			regCenterPrefix: "/conf",
			isErr:           false,
		},
		{
			regCenterPrefix: "/kaixinweiwu",
			isErr:           true,
		},
	}
	Convey("TestconfFile", t, func() {
		So(err, ShouldBeNil)
		So(AimC(), ShouldNotBeNil)
		for _, v := range unitList {
			err := regCenter(v.regCenterPrefix)
			if v.isErr {
				So(err, ShouldBeError)
				t.Log(err)
			} else {

				So(err, ShouldBeNil)

				var conf = &testConf{}
				e := config.C().StructList(conf, "conf", "log")
				So(e, ShouldBeNil)
				t.Log(conf.Level, conf.LogFileDir)
			}
		}

	})
}

func TestLoadConf(t *testing.T) {
	type testConf struct {
		Level      string `json:"level"`
		LogFileDir string `json:"logFileDir"`
	}
	aimConfig := testconf.GetTestConfPath()
	aimEnv := config.EnvLocal
	os.Setenv("AIM_CONFIG", aimConfig)
	os.Setenv("AIM_ENV", aimEnv)

	var unitList = []TestList{
		{
			regCenterPrefix: "/conf",
			confPath:        []string{"srv.yml"},
			isErr:           false,
		},
		{
			regCenterPrefix: "/abc",
			confPath:        []string{"db.yml"},
			isErr:           true,
		},
		{
			regCenterPrefix: "",
			confPath:        []string{},
			isErr:           false,
		},
	}
	Convey("TestLoadConf", t, func() {
		for _, v := range unitList {
			err := LoadConf(v.regCenterPrefix, v.confPath...)
			if v.isErr {
				So(err, ShouldBeError)
				t.Log(err)
			} else {

				So(err, ShouldBeNil)
				So(AimC().Register, ShouldNotBeNil)
				t.Log(AimC().Register)
				var conf = &testConf{}
				e := config.C().StructList(conf, "conf", "log")
				So(e, ShouldBeNil)
				t.Log(conf.Level, conf.LogFileDir)
				So(zaplog.LL(), ShouldNotBeEmpty)
			}
		}

	})
}
