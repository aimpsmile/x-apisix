package upstream

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
)

const ResourceName = "upstreams"

type Upstream struct {
	apiopt *api.Options
}

func (r *Upstream) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Millisecond*time.Duration(conf.MConf().Gateway.Timeout))
}

//	根据ID获取配置信息
func (r *Upstream) findByID(id string) (*OneRsp, error) {
	var b *OneRsp
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

func (r *Upstream) list() (*ListRsp, error) {
	var b *ListRsp
	gwctx, _ := r.getContext()
	req := api.NewRequest(r.apiopt).
		Get().
		Context(gwctx).
		SetHeader("Content-Type", "application/json").
		Resource(ResourceName).
		Retries(1)
	err := req.Do().Into(&b)
	return b, err
}

//	创建
func (r *Upstream) create(data *bytes.Buffer) (string, error) {
	var b *OneRsp
	gwctx, _ := r.getContext()
	err := api.NewRequest(r.apiopt).
		Post().
		Context(gwctx).
		SetHeader("Content-Type", "application/json").
		Resource(ResourceName).
		Retries(conf.MConf().Gateway.Retries).
		Body(data).
		Do().
		Into(&b)
	if err != nil {
		return "", fmt.Errorf("[type]upstream[msg] insert upstreams is error [data]%+v[error_info]%w", data, err)
	}
	return util.MatchID(b.Node.Key), nil
}

//	根据ID更新
func (r *Upstream) updateByID(id string, data *bytes.Buffer) (string, error) {
	var b *OneRsp
	if id == "" {
		return "", fmt.Errorf("[type]upstream[msg]upstream_id not eixst.")
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
		return "", fmt.Errorf("[type]upstream[msg] update upstreams is error [upstreamID]%s[error_info]%w", id, err)
	}
	return id, nil
}

//	根据ID删除
func (r *Upstream) DeleteByID(id string) (string, error) {
	var b *OneRsp
	if id == "" {
		return "", fmt.Errorf("[type]upstream[msg]upstream_id not eixst.")
	}
	gwctx, _ := r.getContext()
	err := api.NewRequest(r.apiopt).
		Delete().
		Context(gwctx).
		SetHeader("Content-Type", "application/json").
		Resource(ResourceName).
		SubResource(id).
		Retries(conf.MConf().Gateway.Retries).
		Do().
		Into(&b)
	if err != nil {
		return id, fmt.Errorf("[type]upstream[msg] delete upstreams is error [upstreamID]%s[error_info]%w", id, err)
	}
	return id, nil
}

//	写入
func (r *Upstream) Write(t *task.TaskMsg, confPath string) (string, error) {
	upstreamID, err := r.FindByDesc(t)
	if err != nil {
		return "", fmt.Errorf("[type]upstream[msg]find desc is error [task]%+v[error_info]%w", t, err)
	}

	data := new(bytes.Buffer)
	confPath = fmt.Sprintf("%s/%s", confPath, t.Svariable.ModuleVeh.UpstreamTpl)
	tplData := gateway.MakeUpstreams(t)
	if err := util.JsonRequestBody(data, confPath, t.Svariable.ModuleVeh.TplFormat, tplData); err != nil {
		return "", fmt.Errorf("[type]upstream[msg]format gateway file is error [tplfile]%s[tpldata]%+v[error_info]%w", confPath, tplData, err)
	} else {
		zaplog.ML().Info("upstream gateway tpl file", zaplog.String("gwtpl_file", confPath))
	}

	if upstreamID != "" {
		return r.updateByID(upstreamID, data)
	} else {
		return r.create(data)
	}
}

//	获取自动生成的服务列表
func (r *Upstream) GetAutoUpstream() (map[string]*Node, error) {
	descMap := make(map[string]*Node)
	list, err := r.list()
	if err != nil {
		err := fmt.Errorf("[type]upstream[msg]get upstream list is error[error_info]%w", err)
		return descMap, err
	}
	for _, v := range list.Node.Nodes {
		if !util.MatchAutoPrefixOfKey(v.Value.Desc) {
			continue
		}
		descMap[v.Value.Desc] = v
	}
	return descMap, nil
}

//	批量清除列表
func (r *Upstream) ClearList() (upstreamIDs []string, errors aimerror.Errors) {
	dMsgs, err := r.GetAutoUpstream()
	if err != nil {
		errors = errors.Add(err)
		return
	}
	for _, v := range dMsgs {

		upstreamID := util.MatchID(v.Key)
		if _, err := r.DeleteByID(upstreamID); err != nil {
			errors = errors.Add(fmt.Errorf("[type]upstream[msg]delete is error[server_id]%s[error_info]%w", upstreamID, err))
		} else {
			upstreamIDs = append(upstreamIDs, upstreamID)
		}
	}
	return
}

//	根据详情获取ID
func (r *Upstream) FindByDesc(t *task.TaskMsg) (string, error) {
	list, err := r.list()
	upstreamID := ""

	if err != nil {
		return upstreamID, err
	}
	desc := util.MakeUpstreamDesc(t.Svariable)
	for _, n := range list.Node.Nodes {
		if n.Value.Desc != desc {
			continue
		}
		upstreamID = n.Key
		break
	}

	return util.MatchID(upstreamID), nil
}

func NewUpstream(apiopt *api.Options) *Upstream {
	return &Upstream{apiopt: apiopt}
}
