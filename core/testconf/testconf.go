package testconf

//	提供基础组件进行单元测试使用
import (
	"fmt"
	"github.com/micro-in-cn/x-apisix/core"
	"path/filepath"

	"github.com/micro-in-cn/x-apisix/core/config"
	"github.com/micro-in-cn/x-apisix/core/lib/file"
)

func GetTestConfPath() string {
	path := filepath.Dir(file.GetRunDir(1) + "/../../")
	return fmt.Sprintf("%s/base/test/config", path)
}

//	初始化环境变量-不加载配置文件
func InitEnv(opts ...MockFunc) (*conf, error) {

	//confDir := fmt.Sprintf("%s/base/test/config",GetTestConfPath())
	c := &conf{
		aimEnv:     config.EnvLocal,
		confDir:    GetTestConfPath(),
		etcdPrefix: "",
		confPath:   []string{"srv.yml", "db.yml", "broker.yml", "cache.yml"},
	}
	//	设置配置选项
	if err := mockOptions(c, opts...); err != nil {
		return nil, err
	}
	var envMap = make(map[string]string)
	//  服务类型
	envMap[config.SnameEnvKey] = "test.passport.srv.v1"
	envMap[config.SverEnvKey] = "v1"
	envMap[config.StypeEnvKey] = "srv"
	envMap[config.SbuEnvKey] = "test"
	envMap[config.SModuleEnvKey] = "passport"
	//	模拟环境是开发、测试、线上
	envMap[config.KeyAimEnv] = c.aimEnv
	//	模拟conf文件路径
	envMap[config.KeyAimConfig] = c.confDir
	for k, v := range envMap {

		if err := setEnv(k, v); err != nil {
			return nil, err
		}
	}
	return c, nil
}

//	初始化全局配置
//	like -> test.MockConf(testconf.AimConfig("/storage/code/aimgo.config"))
func MockConf(opts ...MockFunc) error {
	//	模拟配置与日志初始化
	c, err := InitEnv(opts...)
	if err != nil {
		return err
	}
	err = core.LoadConf(c.etcdPrefix, c.confPath...)
	if err != nil {
		return err
	}
	return nil
}
