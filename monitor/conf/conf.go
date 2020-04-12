package conf

import (
	"github.com/micro-in-cn/x-apisix/core/config"
)

var c = &Conf{Gateway: &Gateway{}, Leader: &Leader{}}

const VerAll = "*"

//	模块版本对应的配置
type ModuleVer struct {
	// route 模板
	RouteTpl string `json:"routeTpl"`
	// upstream 模板
	UpstreamTpl string `json:"upstreamTpl"`
	// service 模板
	ServiceTpl string `json:"serviceTpl"`
	//	模板格式：json、yaml
	TplFormat string `json:"tplFormat"`
	//	域名列表
	Hosts []string `json:"hosts"`
}

//	过滤配置
type Filter struct {
	BU     string                `json:"bu"`
	Stype  string                `json:"stype"`
	Module string                `json:"module"`
	Ver    map[string]*ModuleVer `json:"ver"`
}

//  网关相关的配置
type Gateway struct {
	//  请求接口超时时间毫秒
	Timeout int `json:"timeout"`
	//	请求接口重试次数
	Retries int `json:"retries"`
	//	网关身份认证token
	Apikey string `json:"apikey"`
	//	网关的接口路径
	Baseurl string `json:"baseurl"`
	//	proto文件路径
	ProtoPath string `json:"protoPath"`
	//	禁用同步到网关的路由
	ForbidRoutes map[string]bool `json:"forbidRoutes"`
}

type Leader struct {
	//  选举需要用的id
	ID string `json:"id"`
	//	leader 属组
	Group string `json:"group"`
	//	节点列表
	Nodes []string `json:"nodes"`
}

// 同步检查配置
type Check struct {
	//	请求接口重试次数
	Retries int `json:"retries"`
	//	检查间隔：单位s
	Interval int `json:"interval"`
}

//	gateway配置
type Conf struct {
	//	监控配置选项
	Check *Check `json:"check"`
	//	版本匹配配置
	Filter []*Filter `json:"filter"`
	//	网关的相关的配置
	Gateway *Gateway `json:"gateway"`
	//	leader主从配置
	Leader *Leader `json:"leader"`
}

//	加载配置文件到内存中
func loadConf(confPath string) (*Conf, error) {
	var conf = &Conf{}
	if err := config.C().Append(
		config.FileConf("", confPath),
	); err != nil {
		return nil, err
	}
	if err := config.C().StructList(conf, "conf"); err != nil {
		return nil, err
	}
	return conf, nil
}

//	初始化配置文件
func InitConf(confPath string) (err error) {
	c, err = loadConf(confPath)
	return
}

// 单元测试才可以使用
func MockConf(conf *Conf) {
	c = conf
}

//	获取配置信息
func MConf() *Conf {
	return c
}
