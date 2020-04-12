/*
@Time : 2019/11/24 下午6:35
@Author : songxiuxuan
@File : microlog_test.go
@Software: GoLand
*/
package microlog

import (
	"log"
	"testing"

	"github.com/micro-in-cn/x-apisix/core/testconf"
	"github.com/micro/go-micro/v2/logger"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	if err := testconf.MockConf(); err != nil {
		log.Panic(err)
	}
	l, err := NewLogger(WithNamespace("go-micro"), logger.WithFields(map[string]interface{}{"service": "debug"}))
	if err != nil {
		log.Fatal(err)
	}
	logger.DefaultLogger = l
}
func TestMicroLog(t *testing.T) {
	Convey("TestMicroLog", t, func() {
		logger.Fields(map[string]interface{}{"ceshi": "fields->"})
		logger.Debug("this is micro.debug log .")
		logger.Info("this is micro.info log .")
		logger.Error("this is micro.error log .")
	})
}

func TestName(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	if l.String() != "microlog" {
		t.Errorf("name is error %s", l.String())
	}

	t.Logf("test logger name: %s", l.String())
}

func TestLogf(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	logger.DefaultLogger = l
	logger.Logf(logger.InfoLevel, "test logf: %s", "name")
}
