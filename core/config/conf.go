package config

import (
	"fmt"
	"time"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/config/source/etcd"
	"github.com/micro/go-micro/v2/config/source/file"
	mlog "github.com/micro/go-micro/v2/logger"
)

type WatchFunc func(name ...string) (err error)
type SourceList []source.Source
type Source source.Source

var instance = &configurator{conf: func() config.Config {
	c, _ := config.NewConfig()
	return c
}()}

func FileConf(p, f string) source.Source {

	path := fmt.Sprintf("%s/%s", p, f)
	return file.NewSource(
		file.WithPath(path),
	)
}

//	etcd配置源
func EtcdConf(address, prefix string, timeout int) source.Source {

	var t time.Duration
	if timeout != 0 {
		t = time.Duration(timeout) * time.Second
	} else {
		t = 3 * time.Second
	}
	return etcd.NewSource(
		etcd.WithAddress(address),
		etcd.WithPrefix(prefix),
		etcd.StripPrefix(false),
		etcd.WithDialTimeout(t),
	)
}

type Configurator interface {
	Sync() error
	Watch(watch WatchFunc, path ...string)
	Scan(conf interface{}) error
	MapAll() map[string]interface{}
	SliceList(cList []string, name ...string) (conf []string, err error)
	MapList(cMap map[string]string, name ...string) (conf map[string]string, err error)
	StructList(cStruct interface{}, name ...string) (err error)
	Append(opts ...source.Source) error
}
type configurator struct {
	conf config.Config
}

//	强制更新配置文件
func (c *configurator) Sync() error {
	return c.conf.Sync()
}

//	监听指定的配置文件
func (c *configurator) Watch(watch WatchFunc, path ...string) {
	var watcher config.Watcher
	go func() {

		for i := 0; i <= 100; i++ {
			w, err := c.conf.Watch(path...)
			if err != nil {
				mlog.Error("[CONF_WATCH] conf.Watch() Error:", err)
				continue
			}
			watcher = w
			break
		}

		for {
			// get next
			_, err := watcher.Next()
			if err != nil {
				mlog.Error("[CONF_WATCH] watcher.Next() Error:", err)
				time.Sleep(time.Second)
				continue
			}
			//	配置强制同步
			mlog.Info("[CONF_WATCH] conf.Sync() start ")
			if err := c.conf.Sync(); err != nil {
				mlog.Error("[CONF_WATCH] conf.Sync() return Error:", err)
			}
			//	更新调用方的监控信息
			if err := watch(path...); err != nil {
				mlog.Error("[CONF_WATCH] conf.Watch() return Error:", err)
			}
		}
	}()
}

//	获取所有配置列表
func (c *configurator) Scan(conf interface{}) error {
	return c.conf.Scan(conf)
}

//	获取所有map列表
func (c *configurator) MapAll() map[string]interface{} {
	return c.conf.Map()
}

//	获取slice列表
func (c *configurator) SliceList(cList []string, name ...string) (conf []string, err error) {
	v := c.conf.Get(name...)
	if v != nil {
		conf = v.StringSlice(cList)
	} else {
		err = fmt.Errorf("[SliceList] 配置不存在，err：%+v", name)
	}
	return
}

//	获取map列表
func (c *configurator) MapList(cMap map[string]string, name ...string) (conf map[string]string, err error) {
	v := c.conf.Get(name...)
	if v != nil {
		conf = v.StringMap(cMap)
	} else {
		err = fmt.Errorf("[MapList] 配置不存在，err：%+v", name)
	}
	return
}

// 获取指定的struct结构
func (c *configurator) StructList(cStruct interface{}, name ...string) (err error) {
	v := c.conf.Get(name...)
	if v != nil {
		err = v.Scan(cStruct)
	} else {
		err = fmt.Errorf("[StructList] 配置不存在，err：%s", name)
	}
	return err
}

//	加载配置源
func (c *configurator) loadConfSource(sourceList ...source.Source) (err error) {

	if err := c.conf.Load(sourceList...); err != nil {
		return err
	}

	return
}

//	配置文件，追加配置列表
func (c *configurator) Append(sourceList ...source.Source) error {
	return c.loadConfSource(sourceList...)
}

//	获取配置文件列表
func C() Configurator {
	return instance
}
