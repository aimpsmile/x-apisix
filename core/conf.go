package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/micro-in-cn/x-apisix/core/aimerror"
	"github.com/micro-in-cn/x-apisix/core/config"
	"github.com/micro-in-cn/x-apisix/core/lib/file"
	"github.com/micro-in-cn/x-apisix/core/logger/microlog"
	"github.com/micro-in-cn/x-apisix/core/logger/zaplog"
	"github.com/micro/go-micro/v2/logger"
)

const (
	GlobalConfFile = "global.yml"
	GlobalItem     = "global"
	LogicItem      = "conf"
)

var (
	aimc config.Gconf
	aime = aimerror.Errors([]error{})
)

//	获取启动环境类型
func doAimEnv() (aimEnv string, err error) {
	aimEnv = strings.Trim(os.Getenv("AIM_ENV"), " ")
	aimEnv = strings.ToLower(aimEnv)
	switch aimEnv {
	case config.EnvLocal:
	case config.EnvTest:
	case config.EnvOnline:
	default:
		return "", fmt.Errorf("[%s] env not found ", aimEnv)
	}
	return
}

//	获取本配置路径
func doAimConfig(aimEnv string) (aimConfig string, err error) {
	aimConfig = strings.Trim(os.Getenv("AIM_CONFIG"), " ")
	aimConfig = fmt.Sprintf("%s/%s", strings.TrimRight(aimConfig, "/"), aimEnv)
	//	check file is exists
	if ok, err := file.Exists(aimConfig); !ok {
		if err == nil {
			err = fmt.Errorf("file path [%s] not exists", aimConfig)
		}
		return "", err
	}
	return
}

//	注册配置中心
func regCenter(regCenterPrefix string) error {
	if regCenterPrefix == "" {
		return nil
	}
	var reg config.Source
	switch AimC().Register.Type {
	case config.RegisterEtcd:
		//	etcdctl put /${etcdPrefix}/database '{"host": "10.0.0.2", "port": 6379,"charset":"gb12312","source":"etcd"}'
		reg = config.EtcdConf(AimC().Register.Host, regCenterPrefix, AimC().Register.Timeout)
	default:
		return nil
	}
	if err := config.C().Append(reg); err != nil {
		return err
	}
	if err := config.C().StructList(AimC(), LogicItem); err != nil {
		return err
	}
	return nil
}

//	注册本地配置文件
func confFile(aimConfig string, confPath ...string) error {

	confPath = append(confPath, GlobalConfFile)
	var confList config.SourceList
	for _, f := range confPath {
		confList = append(confList, config.FileConf(aimConfig, f))
	}
	if err := config.C().Append(
		confList...,
	); err != nil {
		return err
	}
	if err := config.C().StructList(AimC(), GlobalItem); err != nil {
		return err
	}

	return nil
}

//	加载基础配置、日志配置
func LoadConf(regCenterPrefix string, confPath ...string) error {
	//	获取启动环境类型
	aimEnv, err := doAimEnv()
	if err != nil {
		return err
	}
	//	获取本配置路径
	aimConfig, err := doAimConfig(aimEnv)
	if err != nil {
		return err
	}
	AimC().EnvType = aimEnv
	AimC().Sname = os.Getenv(config.SnameEnvKey)
	AimC().BU = os.Getenv(config.SbuEnvKey)
	AimC().Module = os.Getenv(config.SModuleEnvKey)
	if os.Getenv(config.MsverEnvKey) != "" {
		AimC().Sver = os.Getenv(config.MsverEnvKey)
	} else {
		AimC().Sver = os.Getenv(config.SverEnvKey)
	}
	AimC().Smetatable = config.ServerMetatable(AimC().Smetatable)
	//	注册本地配置文件
	err = confFile(aimConfig, confPath...)
	if err != nil {
		return err
	}
	if err := regCenter(regCenterPrefix); err != nil {
		return err
	}

	//	初始化micro日志
	err = zaplog.ML().New(AimC().Sname, aimEnv, AimC().Sver)
	if err != nil {
		return err
	}
	//	初始化业务日志
	err = zaplog.LL().New(AimC().Sname, aimEnv, AimC().Sver)
	if err != nil {
		return err
	}
	//	设置micro日志为zap
	l, err := microlog.NewLogger(microlog.WithNamespace("go-micro"))
	if err != nil {
		return err
	}
	logger.DefaultLogger = l
	return nil
}

//	拉取全局配置
func AimC() *config.Gconf {
	return &aimc
}

//	获取错误列表是否
func AimErrors() aimerror.Errors {
	return aime
}

//	只有在逻辑初始化场景才可以使用,,追加错误
func AimErrorsAppend(newErrors ...error) {
	aime = AimErrors().Add(newErrors...)
}
