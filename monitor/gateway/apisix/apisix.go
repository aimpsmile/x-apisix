package apisix

import (
	"fmt"
	"net/http"

	"github.com/micro-in-cn/x-apisix/core/aimerror"
	"github.com/micro-in-cn/x-apisix/core/config"
	"github.com/micro-in-cn/x-apisix/core/lib/api"
	"github.com/micro-in-cn/x-apisix/core/logger/zaplog"
	"github.com/micro-in-cn/x-apisix/monitor/conf"
	"github.com/micro-in-cn/x-apisix/monitor/gateway"
	"github.com/micro-in-cn/x-apisix/monitor/gateway/apisix/proto"
	"github.com/micro-in-cn/x-apisix/monitor/gateway/apisix/route"
	"github.com/micro-in-cn/x-apisix/monitor/gateway/apisix/upstream"
	"github.com/micro-in-cn/x-apisix/monitor/gateway/util"
	"github.com/micro-in-cn/x-apisix/monitor/task"
)

type apisix struct {
	confPath string
	route    *route.Route
	upstream *upstream.Upstream
	proto    *proto.Proto
}

//	action->{insert,update}
func (g *apisix) write(t *task.TaskMsg, isProto bool) (errors aimerror.Errors) {
	//	upstream-写入
	upstreamID, err := g.upstream.Write(t, g.confPath)
	if err != nil {
		errors = errors.Add(err)
		return
	}
	zaplog.ML().Info("[write]upstream successful list.",
		zaplog.String("snmae", t.Svariable.Sname),
		zaplog.String("servicre_id", upstreamID))

	//	proto-写入
	var protoTypes []*util.ProtoType
	if isProto {
		pTypes, perr := g.proto.Write(t, conf.MConf().Gateway.ProtoPath)
		if perr.IsError() {
			errors = errors.Add(perr)
			if len(pTypes) == 0 {
				return
			}
		}
		if len(pTypes) > 0 {
			zaplog.ML().Info("[write]proto successful list.",
				zaplog.String("snmae", t.Svariable.Sname),
				zaplog.String("servicre_id", upstreamID),
				zaplog.Reflect("list", pTypes))
		}
		protoTypes = pTypes
	}

	//	路由-写入
	routeIDs, rerr := g.route.Write(t, isProto, g.confPath, upstreamID, protoTypes)
	if len(routeIDs) > 0 {
		zaplog.ML().Info("[write]route  successful list.",
			zaplog.String("snmae", t.Svariable.Sname),
			zaplog.String("servicre_id", upstreamID),
			zaplog.Strings("list", routeIDs))
	}
	if rerr.IsError() {
		errors = errors.Add(rerr)
	}
	return errors
}

//	action->{delete}
func (g *apisix) delete(t *task.TaskMsg) (errors aimerror.Errors) {
	upstreamID, ferr := g.upstream.FindByDesc(t)
	if ferr != nil {
		zaplog.ML().Info("upstream not exists,not exec delete",
			zaplog.String("sname", t.Svariable.Sname),
			zaplog.NamedError("error_info", ferr))
		return
	}
	routeIDs, rerr := g.route.DeleteByUpstreamID(upstreamID)
	if len(routeIDs) > 0 {
		zaplog.ML().Info("[delete]route successful list.",
			zaplog.String("snmae", t.Svariable.Sname),
			zaplog.String("servicre_id", upstreamID),
			zaplog.Strings("list", routeIDs))
	}
	if rerr.IsError() {
		errors = errors.Add(rerr)
		return
	}
	return errors
}

//	注册中心与网关的route比对
func (g *apisix) routeSameAsGateway(list *route.ListRsp, t *task.TaskMsg) bool {
	forbidRoues := g.route.ForbidRoutes()
	for _, endpoint := range t.Service.Endpoints {
		//	不要同步状态与健康检查
		if _, ok := forbidRoues[endpoint.Name]; ok {
			continue
		}
		desc := util.MakeRouteDesc(t.Svariable, endpoint)
		flag := false
		for _, n := range list.Node.Nodes {
			if !util.MatchAutoPrefixOfKey(n.Value.Desc) {
				continue
			}
			if n.Value.Desc == desc {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

//	注册中心与网关的node比对
func (g *apisix) nodeSameAsGateway(node *upstream.Value, t *task.TaskMsg) bool {
	//	网关与注册中心的upstream相同。需要检查node与route是否相同
	for _, n := range t.Service.Nodes {
		if _, ook := node.Nodes[n.Address]; !ook {
			return false
		}
	}
	return true
}

//	注册中心同步到网关
func (g *apisix) Sync(t *task.TaskMsg) (errors aimerror.Errors) {
	isProto := false
	if t.Svariable.Stype == config.ServiceHTTP2 {
		isProto = true
	}
	switch t.Action {
	case task.ACTION_UPDATE, task.ACTION_CREATE:
		return g.write(t, isProto)
	case task.ACTION_DELETE:
		return g.delete(t)
	}
	return nil
}

//	清除后续操作
func (g *apisix) Cleanup(init bool, dMsgs map[string]string) (delSIDs, delRIDs []string, errors aimerror.Errors) {
	//	clear all proto
	if init {
		zaplog.ML().Info("[start]clear proto.")
		protoIDs, _ := g.proto.ClearList()
		if len(protoIDs) > 0 {
			zaplog.ML().Info("[clear]proto successful list.", zaplog.Strings("list", protoIDs))
		}
		zaplog.ML().Info("[end]clear proto.")
	}

	for _, upstreamID := range dMsgs {
		//	批量删除route
		rIDs, rerr := g.route.DeleteByUpstreamID(upstreamID)
		if len(rIDs) > 0 {
			delRIDs = append(delRIDs, rIDs...)
		}
		if rerr.IsError() {
			errors = errors.Add(rerr)
			continue
		}
		//	批量清理upstream
		_, serr := g.upstream.DeleteByID(upstreamID)
		if serr != nil {
			errors = errors.Add(serr)
		} else {
			delSIDs = append(delSIDs, upstreamID)
		}
	}

	return delSIDs, delRIDs, errors
}

//	task里面的数据，是业务注册中心全部服务数据
func (g *apisix) AllDiff(tMsgs []*task.TaskMsg) (diffTasks []*task.TaskMsg, delMsgs map[string]string, err error) {
	list, rerr := g.route.List()
	if rerr != nil {
		err = fmt.Errorf("[type]route[msg]get route list is error[error_info]%w", rerr)
		return
	}

	delMsgs = make(map[string]string)
	servcieList, serr := g.upstream.GetAutoUpstream()
	if serr != nil {
		err = serr
		return
	}

	for _, t := range tMsgs {
		desc := util.MakeUpstreamDesc(t.Svariable)
		if v, ok := servcieList[desc]; ok {
			//	网关与注册中心相同的服务
			delete(servcieList, desc)
			//	检查node与route是否有差异
			nSame := g.nodeSameAsGateway(v.Value, t)
			rSame := g.routeSameAsGateway(list, t)
			if rSame && nSame {
				continue
			}
		}
		//	注册中心未通步网关或者同步数据有差异
		diffTasks = append(diffTasks, t)
	}
	//	网关大于注册中心
	for k, vv := range servcieList {
		delMsgs[k] = util.MatchID(vv.Key)
	}
	return diffTasks, delMsgs, err
}

func apiOptions() *api.Options {
	c := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	return &api.Options{
		Host:    conf.MConf().Gateway.Baseurl,
		Headers: map[string]string{"X-API-KEY": conf.MConf().Gateway.Apikey},
		Client:  c,
	}
}

func newClient(opts ...gateway.Option) *apisix {
	options := gateway.Options{
		ConfPath: "",
	}
	for _, o := range opts {
		o(&options)
	}
	apiOpts := apiOptions()
	return &apisix{
		route:    route.NewRoute(apiOpts),
		upstream: upstream.NewUpstream(apiOpts),
		proto:    proto.NewProto(apiOpts),
		confPath: options.ConfPath,
	}
}

func NewClient(opts ...gateway.Option) gateway.GatewayI {
	return newClient(opts...)
}
