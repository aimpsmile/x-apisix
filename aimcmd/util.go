package aimcmd

import (
	"strings"
)

//	服务变量
type Svariable struct {
	Sname  string
	BU     string
	Module string
	Stype  string
	Sver   string
}

//	解析micro框架服务名称规则,进而统一整体规范
func ParseSname(sname string) *Svariable {
	service := strings.Split(sname, ".")
	info := &Svariable{Sname: sname}
	if len(service) == 0 {
		return info
	}
	for k, v := range service {
		switch k {
		case 0:
			info.BU = v
		case 1:
			info.Stype = v
		case 2:
			info.Sver = v
		case 3:
			info.Module = v
		}
	}
	return info
}
