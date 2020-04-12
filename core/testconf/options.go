package testconf

//	单元测试配置
import (
	"os"

	apipro "github.com/micro/go-micro/v2/api/proto"
)

//	配置结构体
type conf struct {
	//	环境变量
	aimEnv string
	//	配置目录
	confDir string
	//	etcd前缀
	etcdPrefix string
	//	配置信息
	confPath []string
}

//	模拟API请求结构体
type APIFunc func(req *apipro.Request)

//	模拟全局配置
type MockFunc func(*conf) error

//	设置环境变量
func setEnv(k, v string) error {
	return os.Setenv(k, v)
}

//	加载mock选项
func mockOptions(c *conf, mockFunc ...MockFunc) error {
	for _, o := range mockFunc {
		if err := o(c); err != nil {
			return err
		}
	}
	return nil
}

//	模拟配置etcd前缀
func EtcdPrefix(v string) MockFunc {
	return func(c *conf) error {
		c.etcdPrefix = v
		return nil
	}
}

//	模拟配置文件路径
func ConfPath(v []string) MockFunc {
	return func(c *conf) error {
		c.confPath = v
		return nil
	}
}

//	模拟环境是开发、测试、线上
func AimEnv(v string) MockFunc {
	return func(c *conf) error {
		c.aimEnv = v
		return nil
	}
}

//	模拟conf文件路径
func AimConfig(v string) MockFunc {
	return func(c *conf) error {
		c.confDir = v
		return nil
	}
}

//	模拟micro的环境变量
func MicroEnv(k, v string) MockFunc {
	return func(c *conf) error {
		return setEnv(k, v)
	}
}
