package config

//	执行测试命令
//  go test -v -count=1  *.go
import (
	"fmt"
	"github.com/micro-in-cn/x-apisix/core/lib/file"
	"github.com/smartystreets/goconvey/convey"
	"log"
	"path/filepath"
	"testing"
	"time"
)

type Database struct {
	//	使用这种方式传递 "tcp(host:port)"
	Host string `json:"host"`
	// 用户名.
	User string `json:"user"`
	//	密码
	Password string `json:"password"`
	//	数据库名称
	Db string `json:"db"`
	//	字符集
	Charset string `json:"charset"`
	// 最大打开连接数
	MaxOpenConns int `json:"maxOpenConns"`
	// 保留最大的空闲连接
	MaxIdleConns int `json:"maxIdleConns"`
	// 连接最大生存时间
	ConnMaxLifetime time.Duration `json:"connMaxLifetime"`
	// true:打开日志、false:关闭日志
	LogMode bool `json:"logMode"`
}
type testConf struct {
	Level      string `json:"level"`
	LogFileDir string `json:"logFileDir"`
}

func getConf() string {
	return filepath.Dir(file.GetRunDir(1) + "/../../")
}

func initConf() {
	confDir := fmt.Sprintf("%s/base/test/config", getConf())
	aimEnv := EnvLocal
	aimConfig := fmt.Sprintf("%s/%s", confDir, aimEnv)
	confPath := []string{"global.yml", "srv.yml", "db.yml"}

	var confList SourceList
	for _, f := range confPath {
		confList = append(confList, FileConf(aimConfig, f))
	}
	if err := C().Append(
		confList...,
	); err != nil {
		log.Panic(err)
	}

}

func TestInit(t *testing.T) {
	convey.Convey("TestEtcdConf", t, func() {
		conf := C().MapAll()
		t.Log(conf)
	})
}

func TestEtcdConf(t *testing.T) {
	convey.Convey("TestEtcdConf", t, func() {
		initConf()
		var (
			aimc = Gconf{}
			conf = &testConf{}
		)
		//	加载全局配置
		err := C().StructList(&aimc, "global")
		convey.So(err, convey.ShouldBeNil)
		//	追踪注册中心
		err = C().Append(EtcdConf(aimc.Register.Host, "/conf", aimc.Register.Timeout))
		convey.So(err, convey.ShouldBeNil)

		//	从注册中心获取配置数据
		e := C().StructList(conf, "conf", "log")
		convey.So(e, convey.ShouldBeNil)
		t.Log(conf.Level, conf.LogFileDir)
	})

}

func TestStructList(t *testing.T) {
	var config Database
	convey.Convey("TestStructList", t, func() {
		initConf()
		err := C().StructList(&config, "db", "mysql", "aim_uqudu_db")
		convey.So(err, convey.ShouldBeNil)
		convey.So(config.Host, convey.ShouldNotBeEmpty)
	})
}

func TestMapList(t *testing.T) {
	convey.Convey("TestMapList", t, func() {
		initConf()
		conf, err := C().MapList(map[string]string{}, "db", "mysql", "aim_uqudu_db")
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(conf), convey.ShouldBeGreaterThan, 0)
	})
}

func TestMapAll(t *testing.T) {
	convey.Convey("TestMapAll", t, func() {
		initConf()
		conf := C().MapAll()
		convey.So(len(conf), convey.ShouldBeGreaterThan, 0)
	})
}

func TestSliceList(t *testing.T) {
	convey.Convey("TestSliceList", t, func() {
		initConf()
		conf, err := C().SliceList([]string{"weiwu", "hahah"}, "conf", "depend", "db", "aim_uqudu_db")
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(conf), convey.ShouldBeGreaterThan, 0)
		t.Log(conf)
	})
}

func TestServiceTypeList(t *testing.T) {
	t.Log(ServiceTypeList)
}

func BenchmarkConf(b *testing.B) {

	// 灾难恢复
	defer func() {
		if r := recover(); r != nil {
			b.Log("[main] Recovered in f %V", r)
		}
	}()
	confDir := fmt.Sprintf("%s/base/test/config", getConf())
	aimEnv := EnvLocal
	aimConfig := fmt.Sprintf("%s/%s", confDir, aimEnv)
	confPath := []string{"global.yml", "srv.yml", "db.yml"}

	var confList SourceList
	for _, f := range confPath {
		confList = append(confList, FileConf(aimConfig, f))
	}
	err := C().Append(
		confList...,
	)
	if err != nil {
		b.Error(err)
	}
	var config Database
	for i := 0; i <= b.N; i++ {
		err := C().StructList(&config, "db", "mysql", "aim_uqudu_db")
		if err != nil {
			b.Error(err)
		}
	}
}
