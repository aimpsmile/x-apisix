package route

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/micro-in-cn/x-apisix/core/aimerror"
	"github.com/micro-in-cn/x-apisix/core/lib/api"
	"github.com/micro-in-cn/x-apisix/core/logger/zaplog"
	"github.com/micro-in-cn/x-apisix/monitor/conf"
	"github.com/micro-in-cn/x-apisix/monitor/gateway"
	"github.com/micro-in-cn/x-apisix/monitor/gateway/util"
	"github.com/micro-in-cn/x-apisix/monitor/task"
	"github.com/micro/go-micro/v2/registry"
)

const ResourceName = "routes"

type Route struct {
	apiopt *api.Options
}

var NotMatchRoue = map[string]bool{
	//	服务的状态
	"/stats": true,
	//	健康检查
	"/health": true,
}

//	禁用的路由列表
func (r *Route) ForbidRoutes() (forbidRoues map[string]bool) {
	forbidRoues = conf.MConf().Gateway.ForbidRoutes
	if forbidRoues == nil {
		forbidRoues = make(map[string]bool)
	}
	return
}
func (r *Route) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Millisecond*time.Duration(conf.MConf().Gateway.Timeout))
}

//	根据ID获取配置信息
func (r *Route) findByID(id string) (*OneRsp, error) {
	var b *OneRsp
	if id == "" {
		return nil, fmt.Errorf("[type]route[msg]route_id not eixst.")
	}
	gwctx, _ := r.getContext()
	err := api.NewRequest(r.apiopt).
		Get().
		Context(gwctx).
		SetHeader("Content-Type", "application/json").
		Resource(ResourceName).
		SubResource(id).
		Retries(conf.MConf().Gateway.Retries).
		Do().
		Into(&b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *Route) List() (*ListRsp, error) {
	var b *ListRsp
	gwctx, _ := r.getContext()
	req := api.NewRequest(r.apiopt).
		Context(gwctx).
		Get().
		SetHeader("Content-Type", "application/json").
		Resource(ResourceName).
		Retries(1)
	err := req.Do().Into(&b)
	return b, err
}

//	匹配protoID
func (r *Route) matchProtoID(packageName string, endpoint *registry.Endpoint, pList []*util.ProtoType) string {
	for _, pval := range pList {
		if packageName != pval.PackageName {
			continue
		}
		for _, sname := range pval.ServiceList {
			if endpoint.Metadata["grpc_service"] == sname {
				return pval.ProtoID
			}
		}
	}
	return ""
}

//	根据service_id获取路由配置信息
func (r *Route) findByServiceID(serviceID string) ([]*Node, error) {
	if serviceID == "" {
		return nil, fmt.Errorf("[type]route[msg]service_id not eixst.")
	}
	list, err := r.List()
	var rsp []*Node

	if err != nil {
		return nil, err
	}
	for _, n := range list.Node.Nodes {
		if n.Value.ServiceID != serviceID {
			continue
		}
		rsp = append(rsp, n)
	}

	return rsp, nil
}

//	根据upstream_id获取路由配置信息
func (r *Route) findByUpstreamID(upstreamID string) ([]*Node, error) {
	if upstreamID == "" {
		return nil, fmt.Errorf("[type]route[msg]upstream_id not eixst.")
	}
	list, err := r.List()
	var rsp []*Node

	if err != nil {
		return nil, err
	}
	for _, n := range list.Node.Nodes {
		if n.Value.UpstreamID != upstreamID {
			continue
		}
		rsp = append(rsp, n)
	}

	return rsp, nil
}

//	创建
func (r *Route) create(data *bytes.Buffer) (string, error) {
	var b *OneRsp
	gwctx, _ := r.getContext()
	err := api.NewRequest(r.apiopt).
		Post().
		Context(gwctx).
		SetHeader("Content-Type", "application/json").
		Resource(ResourceName).
		Body(data).
		Do().
		Into(&b)
	if err != nil {
		return "", fmt.Errorf("[type]route[msg] insert route is error [data]%+v[error_info]%w", data, err)
	}
	return util.MatchID(b.Node.Key), nil
}

//	根据ID更新
func (r *Route) updateByID(id string, data *bytes.Buffer) (string, error) {
	var b *OneRsp
	if id == "" {
		return "", fmt.Errorf("[type]route[msg]route_id not eixst.")
	}
	gwctx, _ := r.getContext()
	err := api.NewRequest(r.apiopt).
		Put().
		Context(gwctx).
		SetHeader("Content-Type", "application/json").
		Resource(ResourceName).
		SubResource(id).
		Retries(conf.MConf().Gateway.Retries).
		Body(data).
		Do().
		Into(&b)
	if err != nil {
		return "", fmt.Errorf("[type]route[msg] update route is error [routeID]%s[error_info]%w", id, err)
	}
	return util.MatchID(b.Node.Key), nil
}

//	写入数据
func (r *Route) write(t *task.TaskMsg, endpoint *registry.Endpoint, confPath string, tplData *gateway.Routes) (string, error) {
	routeID, err := r.FindByDesc(t, endpoint)
	if err != nil {
		return "", fmt.Errorf("[type]route[msg]find desc is error [task]%+v[error_info]%w", t, err)
	}
	data := new(bytes.Buffer)
	if err := util.JsonRequestBody(data, confPath, t.Svariable.ModuleVeh.TplFormat, tplData); err != nil {
		return "", fmt.Errorf("[type]route[msg]format  file is error [tplfile]%s[tpldata]%+v[error_info]%w", confPath, tplData, err)
	} else {
		zaplog.ML().Info("upstream gateway tpl file", zaplog.String("gwtpl_file", confPath))
	}
	if routeID != "" {
		return r.updateByID(routeID, data)
	} else {
		return r.create(data)
	}
}

//	根据ID删除
func (r *Route) DeleteByID(id string) (string, error) {
	var b *OneRsp
	if id == "" {
		return "", fmt.Errorf("[type]route[msg]route_id not eixst.")
	}
	gwctx, _ := r.getContext()
	err := api.NewRequest(r.apiopt).
		Delete().
		Context(gwctx).
		SetHeader("Content-Type", "application/json").
		Resource(ResourceName).
		Retries(conf.MConf().Gateway.Retries).
		SubResource(id).
		Do().
		Into(&b)
	if err != nil {
		return "", fmt.Errorf("[type]route[msg] delete route is error [upstreamID]%s[error_info]%w", id, err)
	}
	return id, nil
}

// 根据service_id删除路由
// nolint dupl
func (r *Route) DeleteByServiceID(serviceID string) (routeIDs []string, errors aimerror.Errors) {
	if serviceID == "" {
		errors = errors.Add(fmt.Errorf("[type]route[msg]service_id not eixst."))
		return
	}
	nodes, err := r.findByServiceID(serviceID)
	if err != nil {
		errors = errors.Add(err)
		return
	}
	for _, v := range nodes {
		id, err := r.DeleteByID(util.MatchID(v.Key))
		if err != nil {
			errors = errors.Add(err)
		} else {
			routeIDs = append(routeIDs, id)
		}
	}
	return
}

//	根据upstream_id删除路由
// nolint dupl
func (r *Route) DeleteByUpstreamID(upstreamID string) (routeIDs []string, errors aimerror.Errors) {
	if upstreamID == "" {
		errors = errors.Add(fmt.Errorf("[type]route[msg]upstream_id not eixst."))
		return
	}
	nodes, err := r.findByUpstreamID(upstreamID)
	if err != nil {
		errors = errors.Add(err)
		return
	}
	for _, v := range nodes {
		id, err := r.DeleteByID(util.MatchID(v.Key))
		if err != nil {
			errors = errors.Add(err)
		} else {
			routeIDs = append(routeIDs, id)
		}
	}
	return
}

//	批量写入
func (r *Route) Write(t *task.TaskMsg, isProto bool, confPath,
	upstreamID string, pList []*util.ProtoType) (routeIDs []string, errors aimerror.Errors) {
	if upstreamID == "" {
		errors = errors.Add(fmt.Errorf("upstreamID not exists."))
		return
	}
	confPath = fmt.Sprintf("%s/%s", confPath, t.Svariable.ModuleVeh.RouteTpl)
	protoPackage := t.Svariable.ProtoPackage
	forbidRoues := r.ForbidRoutes()
	for _, endpoint := range t.Service.Endpoints {
		//	不要同步状态与健康检查
		if _, ok := forbidRoues[endpoint.Name]; ok {
			continue
		}

		//  处理proto生成逻辑
		protoID := ""
		if isProto {
			protoID = r.matchProtoID(protoPackage, endpoint, pList)
			if protoID == "" {
				errors = errors.Add(fmt.Errorf("[type]route[msg]proto_id not exists[proto.package]%s[endpoint]%s",
					protoPackage,
					endpoint.Metadata["grpc_method"]))
				continue
			}
		}
		//  写入
		tplData := gateway.MakeRoutes(t, endpoint, upstreamID, protoID)
		id, err := r.write(t, endpoint, confPath, tplData)
		if err != nil {
			errors = errors.Add(err)
		} else {
			routeIDs = append(routeIDs, id)
		}
	}
	return routeIDs, errors
}

//	根据详情获取ID
func (r *Route) FindByDesc(t *task.TaskMsg, endpoint *registry.Endpoint) (string, error) {
	list, err := r.List()
	routeID := ""

	if err != nil {
		return routeID, err
	}
	desc := util.MakeRouteDesc(t.Svariable, endpoint)
	for _, n := range list.Node.Nodes {
		if n.Value.Desc != desc {
			continue
		}
		routeID = n.Key
		break
	}

	return util.MatchID(routeID), nil
}

//	返回gateway需要的interface接口
func NewRoute(apiopt *api.Options) *Route {
	return &Route{apiopt: apiopt}
}
