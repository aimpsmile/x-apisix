package upstream

import (
	"bytes"
	"github.com/micro-in-cn/x-apisix/core/lib/api"
	"github.com/micro-in-cn/x-apisix/monitor/conf"
	"github.com/micro-in-cn/x-apisix/monitor/gateway"
	"github.com/micro-in-cn/x-apisix/monitor/gateway/util"
	"github.com/micro-in-cn/x-apisix/monitor/task"
	"github.com/micro-in-cn/x-apisix/monitor/test"
	"github.com/micro/go-micro/v2/registry"
	"testing"
)

const confPath = "/storage/code/aimgo.config/local/apisix"

func aTestUpstream(t *testing.T) {
	tpl := `
{{- $len := sub (len .Nodes) 1 }}
{
	"desc": "{{ .Desc }}",
	"key":"",
	"type": "roundrobin",
	"nodes": {
		{{- range $key,$val := .Nodes }}
		"{{ $val.Address }}": 1{{- if lt $key $len }},{{- end }}
		{{- end }}
	}
}
`
	nodes := []*registry.Node{
		{
			Id:       "aimgo.passport.http2.v1-bdd733c8-d773-4b23-a06e-84d307203878",
			Address:  "127.0.0.1:57628",
			Metadata: nil,
		},
		{
			Id:       "aimgo.passport.http2.v1-bdd733c8-d773-4b23-a06e-84d307203878",
			Address:  "127.0.0.1:57629",
			Metadata: nil,
		},
	}
	svariable := &task.Svariable{
		Sname:  "aimgo.passport.http2.v1",
		BU:     "aimgo",
		Module: "passport",
		Stype:  "http2",
		Sver:   "v1",
		ModuleVeh: &conf.ModuleVer{
			RouteTpl:    "grpc.routes.json",
			ServiceTpl:  "grpc.services.json",
			UpstreamTpl: "grpc.upstreams.json",
			Hosts:       []string{"srv.uqudu.com", "grpc.uqudu.com"},
		},
	}

	data := gateway.Servcies{
		Desc:      util.MakeUpstreamDesc(svariable),
		Nodes:     nodes,
		Svariable: svariable,
	}
	b := new(bytes.Buffer)
	// 数据驱动模板
	err := util.FormatGateway(b, "upstreams", tpl, data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)
	reqBody := &test.Testcase{
		ReqFn: func(opts *api.Options) *api.Request {
			return api.NewRequest(opts).Post().SetHeader("Content-Type", "application/json").Resource("upstream").Body(b)
		},
		Method: "POST",
		URI:    "/upstream",
		Body:   b.Bytes(),
	}
	headers := map[string]string{"X-API-KEY": conf.MConf().Gateway.Apikey}
	if err := test.Request(reqBody, headers, t); err != nil {
		t.Fatal(err)
	}
}

func aTestFindByDesc(t *testing.T) {
	nodes := []*registry.Node{
		{
			Id:       "aimgo.passport.http2.v1-bdd733c8-d773-4b23-a06e-84d307203878",
			Address:  "127.0.0.1:57628",
			Metadata: nil,
		},
		{
			Id:       "aimgo.passport.http2.v1-bdd733c8-d773-4b23-a06e-84d307203878",
			Address:  "127.0.0.1:57629",
			Metadata: nil,
		},
	}
	svariable := &task.Svariable{
		Sname:  "aimgo.passport.http2.v1",
		BU:     "aimgo",
		Module: "passport",
		Stype:  "http2",
		Sver:   "v1",
		ModuleVeh: &conf.ModuleVer{
			RouteTpl:   "grpc.routes.json",
			ServiceTpl: "grpc.upstreams.json",
			Hosts:      []string{"srv.uqudu.com", "grpc.uqudu.com"},
		},
	}
	data := &task.TaskMsg{
		Action: task.ACTION_CREATE,
		Service: &registry.Service{
			Name:  svariable.Sname,
			Nodes: nodes,
		},
		Svariable: svariable,
	}

	s := NewUpstream(test.HttpOptions())
	//	根据详情获取upstream_id
	upstreamID, err := s.FindByDesc(data)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("findServiceIDByDesc[upstream_id]", upstreamID)
	}
	//	写入操作
	sID, serr := s.Write(data, confPath)
	if serr != nil {
		t.Fatal(serr)
	} else {
		t.Log("update upstream success [upstream_id]", sID)
	}
	//	根据upstream_id获取配置信息
	ss := NewUpstream(test.HttpOptions())
	onersp, ferr := ss.findByID(sID)
	if ferr != nil {
		t.Fatal(ferr)
	} else {
		t.Log(onersp.Node.Key, onersp.Node.Value.Desc)
	}
	//	写入操作
	sid, derr := s.DeleteByID(sID)
	if derr != nil {
		t.Fatal(derr)
	} else {
		t.Log("delete upstream success [upstream_id]", sid)
	}
}

func TestUpstreamYaml(t *testing.T) {
	nodes := []*registry.Node{
		{
			Id:       "aimgo.passport.http2.v1-bdd733c8-d773-4b23-a06e-84d307203878",
			Address:  "127.0.0.1:57628",
			Metadata: nil,
		},
		{
			Id:       "aimgo.passport.http2.v1-bdd733c8-d773-4b23-a06e-84d307203878",
			Address:  "127.0.0.1:57629",
			Metadata: nil,
		},
	}
	svariable := &task.Svariable{
		Sname:  "aimgo.passport.http2.v1",
		BU:     "aimgo",
		Module: "passport",
		Stype:  "http2",
		Sver:   "v1",
		ModuleVeh: &conf.ModuleVer{
			RouteTpl:   "grpc.routes.json",
			ServiceTpl: "grpc.upstreams.json",
			Hosts:      []string{"srv.uqudu.com", "grpc.uqudu.com"},
		},
	}
	data := gateway.Servcies{
		Desc:      util.MakeUpstreamDesc(svariable),
		Nodes:     nodes,
		Svariable: svariable,
	}

	b := new(bytes.Buffer)
	// 数据驱动模板
	err := util.JsonRequestBody(b, "/storage/code/aimgo.config/local/apisix/yaml/web.upstreams.yaml", "yaml", data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)
}
