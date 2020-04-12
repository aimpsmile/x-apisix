package proto

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/micro-in-cn/x-apisix/core/aimerror"
	"github.com/micro-in-cn/x-apisix/core/lib/api"
	"github.com/micro-in-cn/x-apisix/monitor/conf"
	"github.com/micro-in-cn/x-apisix/monitor/gateway/util"
	"github.com/micro-in-cn/x-apisix/monitor/task"
	jsoniter "github.com/json-iterator/go"
)

const ResourceName = "proto"

type Proto struct {
	apiopt *api.Options
}

func (r *Proto) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Millisecond*time.Duration(conf.MConf().Gateway.Timeout))
}

//	根据ID获取配置信息
func (r *Proto) findByID(id string) (*OneRsp, error) {
	var b *OneRsp
	if id == "" {
		return nil, fmt.Errorf("[type]proto[msg]proto_id not eixst.")
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

func (r *Proto) list() (*ListRsp, error) {
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

//	更新
func (r *Proto) update(id string, data *bytes.Buffer) (string, error) {
	var b *OneRsp
	if id == "" {
		return "", fmt.Errorf("[type]proto[msg]proto_id not eixst.")
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
		return "", err
	}
	return util.MatchID(b.Node.Key), nil
}

//	根据ID删除
func (r *Proto) DeleteByID(id string) error {
	if id == "" {
		return fmt.Errorf("[type]proto[msg]proto_id not eixst.")
	}
	gwctx, _ := r.getContext()
	return api.NewRequest(r.apiopt).
		Delete().
		Context(gwctx).
		SetHeader("Content-Type", "application/json").
		Resource(ResourceName).
		Retries(conf.MConf().Gateway.Retries).
		SubResource(id).
		Do().Error()
}

//	批量清除列表
func (r *Proto) ClearList() (protoIDs []string, errors aimerror.Errors) {

	list, err := r.list()
	if err != nil {
		errors = errors.Add(fmt.Errorf("[type]proto[msg]get proto list is error[error_info]%w", err))
		return
	}
	for _, v := range list.Node.Nodes {
		protoID := util.MatchID(v.Key)
		if err := r.DeleteByID(protoID); err != nil {
			errors = errors.Add(fmt.Errorf("[type]proto[msg]delete is error[proto_id]%s[error_info]%w", protoID, err))
		} else {
			protoIDs = append(protoIDs, protoID)
		}
	}
	return
}

//	写入
func (r *Proto) Write(t *task.TaskMsg, protoPath string) (pList []*util.ProtoType, errors aimerror.Errors) {
	files, err := util.GetProtoFileList(protoPath, t.Svariable)
	if err != nil {
		return nil, errors.Add(fmt.Errorf("[type]proto[msg]glob *Proto file is error[proto_path]%s[error_info]%w", protoPath, err))
	}
	for _, f := range files {

		c, err := ioutil.ReadFile(f)
		if err != nil {
			errors = errors.Add(fmt.Errorf("[type]proto[msg]read proto file is error[proto_file]%s[error_info]%w", f, err))
			continue
		}
		protot, perr := util.MakeProtoID(c)
		if perr != nil {
			errors = errors.Add(perr)
			continue
		}
		protoContent, cerr := util.CompressProto(c)
		if cerr != nil {
			errors = errors.Add(cerr)
			continue
		}

		val := &Value{Content: protoContent}
		b := new(bytes.Buffer)
		if jerr := jsoniter.NewEncoder(b).Encode(val); jerr != nil {
			errors = errors.Add(fmt.Errorf("[type]proto[msg]make json is error[protot]%+v[error_info]%w", protot, jerr))
			continue
		}
		if _, uerr := r.update(protot.ProtoID, b); uerr != nil {
			errors = errors.Add(fmt.Errorf("[type]proto[msg] call apisix api is error [protot]%+v[error_info]%w", protot, uerr))
		} else {
			pList = append(pList, protot)
		}
	}
	return pList, errors
}

func NewProto(apiopt *api.Options) *Proto {
	return &Proto{apiopt: apiopt}
}
