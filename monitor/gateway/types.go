package gateway

import (
	"strings"

	"github.com/micro-in-cn/x-apisix/monitor/gateway/util"
	"github.com/micro-in-cn/x-apisix/monitor/task"
	"github.com/micro/go-micro/v2/registry"
)

//	路由需要的结构体
type Routes struct {
	Desc      string
	ProtoID   string
	ServiceID string
	Endpoint  *registry.Endpoint
	Svariable *task.Svariable
	Methods   []string
}

//	服务需要的结构体
type Servcies struct {
	Desc      string
	Nodes     []*registry.Node
	Svariable *task.Svariable
}

type Result struct {
	Key string
	Val string
	Err error
}

//	routes-生成格式化变量
func MakeRoutes(t *task.TaskMsg, endpoint *registry.Endpoint, serviceID, protoID string) *Routes {
	methods := []string{}
	if endpoint.Metadata != nil {
		methods = strings.Split(endpoint.Metadata["method"], ",")
	}
	return &Routes{
		Desc:      util.MakeRouteDesc(t.Svariable, endpoint),
		ServiceID: serviceID,
		ProtoID:   protoID,
		Endpoint:  endpoint,
		Svariable: t.Svariable,
		Methods:   methods,
	}
}

//	servcies-生成格式化变量
func MakeServcies(t *task.TaskMsg) *Servcies {
	return &Servcies{
		Desc:      util.MakeServiceDesc(t.Svariable),
		Nodes:     t.Service.Nodes,
		Svariable: t.Svariable,
	}
}

//	upstream-生成格式化变量
func MakeUpstreams(t *task.TaskMsg) *Servcies {
	return &Servcies{
		Desc:      util.MakeUpstreamDesc(t.Svariable),
		Nodes:     t.Service.Nodes,
		Svariable: t.Svariable,
	}
}

//	网关统一的返回结构体
func MakeResult(key, val string, err error) *Result {
	return &Result{
		Key: key,
		Val: val,
		Err: err,
	}
}

//	网关统一的返回结构体列表
func MakeResults(key, val string, err error) []*Result {
	return []*Result{MakeResult(key, val, err)}
}
