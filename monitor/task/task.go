package task

import (
	"fmt"

	"github.com/micro-in-cn/x-apisix/aimcmd"
	"github.com/micro-in-cn/x-apisix/monitor/conf"
	"github.com/micro/go-micro/v2/registry"
)

const (
	ACTION_CREATE = "create"
	ACTION_UPDATE = "update"
	ACTION_DELETE = "delete"
)

//	服务变量
type Svariable struct {
	Sname        string
	BU           string
	Module       string
	Stype        string
	Sver         string
	Version      string
	ProtoPackage string

	ModuleVeh *conf.ModuleVer
}

//	任务消息体
type TaskMsg struct {
	Action    string
	Retries   int
	Service   *registry.Service
	Svariable *Svariable
}

//	生成任务的信息列表
func NewMsg(action string, service *registry.Service) (*TaskMsg, error) {
	svar, err := ShouldSnameBeConf(service.Name, service.Version, conf.MConf())
	if err != nil {
		return nil, err
	}
	return &TaskMsg{
		Action:    action,
		Retries:   0,
		Service:   service,
		Svariable: svar,
	}, nil
}

//	服务名称，解析成结构体
func ServiceNameToInfo(sname, version string) *Svariable {

	info := aimcmd.ParseSname(sname)
	//  生成proto.package名称
	protoPackge := fmt.Sprintf("%s.%s.%s.%s", info.BU, info.Module, info.Stype, info.Sver)
	return &Svariable{
		Sname:        info.Sname,
		BU:           info.BU,
		Module:       info.Module,
		Stype:        info.Stype,
		Sver:         info.Sver,
		Version:      version,
		ProtoPackage: protoPackge,
		ModuleVeh:    &conf.ModuleVer{},
	}
}

//	版本信息与配置文件进行匹配
func matchVer(version string, filter *conf.Filter) *conf.ModuleVer {

	va, err := StringConverFloat64(version, ".", 64)

	if err == nil {
		for ver, verConf := range filter.Ver {
			vbs, op, vcs := StringConverOperator(ver)
			vb, berr := StringConverFloat64(vbs, ".", 64)
			if berr != nil {
				continue
			}
			vc, cerr := StringConverFloat64(vcs, ".", 64)
			if cerr != nil {
				continue
			}
			if VersionCompare(va, vb, op, vc) {
				return verConf
			}
		}
	}
	if v, ok := filter.Ver[conf.VerAll]; ok {
		return v
	}
	return nil
}

//	服务名称应该配置好
func ShouldSnameBeConf(sname, version string, c *conf.Conf) (*Svariable, error) {
	info := ServiceNameToInfo(sname, version)
	if info.Sname == "" || info.BU == "" ||
		info.Module == "" || info.Stype == "" || info.Sver == "" {
		estr := "[info]sname format not match[sname]%s[-][version]%s[BU]%s[Module]%s[Stype]%s[Sver]%s"
		err := fmt.Errorf(estr,
			sname, version, info.BU, info.Module, info.Stype, info.Sver)
		return nil, err
	}
	for _, filter := range c.Filter {
		if info.BU != filter.BU || info.Stype != filter.Stype {
			continue
		}
		if filter.Module == info.Module {
			info.ModuleVeh = matchVer(info.Version, filter)
			break
		}
		if filter.Module == conf.VerAll {
			info.ModuleVeh = matchVer(info.Version, filter)
		}
	}
	if (info.ModuleVeh.UpstreamTpl != "" || info.ModuleVeh.ServiceTpl != "") &&
		info.ModuleVeh.RouteTpl != "" && len(info.ModuleVeh.Hosts) != 0 {
		return info, nil
	}
	estr := "[info]not conf to match[sname]%s[version]%s[-][BU]%s[Module]%s[Stype]%s[Sver]%s[-][UpstreamTpl]%s[ServiceTpl]%s[RouteTpl]%s[Hosts]%s"
	err := fmt.Errorf(estr,
		sname, version, info.BU, info.Module, info.Stype, info.Sver, info.ModuleVeh.UpstreamTpl,
		info.ModuleVeh.ServiceTpl, info.ModuleVeh.RouteTpl, info.ModuleVeh.Hosts)
	return nil, err
}
