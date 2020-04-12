package route

import (
	"bytes"
	"fmt"
	"github.com/micro-in-cn/x-apisix/monitor/conf"
	"github.com/micro-in-cn/x-apisix/monitor/gateway/util"
	"github.com/micro-in-cn/x-apisix/monitor/test"
	"testing"

	"github.com/micro-in-cn/x-apisix/core/lib/api"
	"github.com/micro-in-cn/x-apisix/monitor/gateway"
	"github.com/micro-in-cn/x-apisix/monitor/task"
	"github.com/micro/go-micro/v2/registry"
)

const confPath = "/storage/code/aimgo.config/local/apisix"

func aTestHttpRoutes(t *testing.T) {
	tpl := `
{{- $lenHosts := sub (len .Svariable.ModuleVeh.Hosts) 1 }}
{
	"desc": "{{ .Desc }}",
	"priority":"9",
	"methods": [
		"POST","GET","PUT","DELETE","PATCH"
	],
	"uris": [
		"/{{ .Svariable.Sver }}{{ .Svariable.Module }}/{{ trimSuffix .Endpoint.Name "/" }}*"
	],
	"hosts": [
		{{- range $key,$val := .Svariable.ModuleVeh.Hosts }}
				"{{- $val }}"{{- if lt $key $lenHosts }},{{- end }}
		{{- end }}
	],
	"service_protocol": "http",
	"service_id": "{{ .ServiceID }}",
	"plugins": {}
}
`
	endpoint := &registry.Endpoint{
		Name: "/",
	}
	svariable := &task.Svariable{
		Sname:  "aimgo.passport.web.v1",
		BU:     "aimgo",
		Module: "passport",
		Stype:  "web",
		Sver:   "v1",
		ModuleVeh: &conf.ModuleVer{
			RouteTpl:    "web.routes.json",
			ServiceTpl:  "web.services.json",
			UpstreamTpl: "grpc.upstreams.json",
			Hosts:       []string{"web.uqudu.com", "web.uqudu.com"},
		},
	}
	data := gateway.Routes{

		Desc:      util.MakeRouteDesc(svariable, endpoint),
		ServiceID: "000022223",
		Endpoint:  endpoint,
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
			return api.NewRequest(opts).Post().SetHeader("Content-Type", "application/json").Resource("route").Body(b)
		},
		Method: "POST",
		URI:    "/route",
		Body:   b.Bytes(),
	}
	headers := map[string]string{"X-API-KEY": conf.MConf().Gateway.Apikey}
	if err := test.Request(reqBody, headers, t); err != nil {
		t.Fatal(err)
	}

}

func aTestGrpcRoutes(t *testing.T) {
	tpl := `
{{- $methods := split .Endpoint.Metadata.method "," }}
{{- $lenMethod := sub (len $methods) 1 }}
{{- $lenHosts := sub (len .Svariable.ModuleVeh.Hosts) 1 }}
{
	"desc": "{{ .Desc }}",
	"priority":"10",
	"methods": [
		{{- range $k,$v := $methods }}
			"{{ $v }}"{{- if lt $k $lenMethod }},{{- end }}
		{{- end }}
	],
	"uris": [
		"/{{ .Svariable.Sver }}{{ .Svariable.Module }}/{{ trimPrefix .Endpoint.Metadata.path "/" }}"
	],
	"hosts": [
		{{- range $key,$val := .Svariable.ModuleVeh.Hosts }}
				"{{ $val }}"{{- if lt $key $lenHosts }},{{- end }}
		{{- end }}
	],
	"plugins": {
		"grpc-transcode": {
				"proto_id": "{{ .ProtoID }}",
				"service": "{{ .Svariable.Sname }}.{{ .Endpoint.Metadata.grpc_service }}",
				"method": "{{ .Endpoint.Metadata.grpc_method }}",
				"pb_option":["int64_as_string"]
		},
		"gheader-plugin" : {}
	},
	"service_protocol": "grpc",
	"service_id": "{{ .ServiceID }}"
}
`
	endpoint := &registry.Endpoint{
		Name: "SayService.Hello",
		Metadata: map[string]string{
			"grpc_service": "SayService",
			"grpc_method":  "Hello",
			"handler":      "HTTP2",
			"method":       "GET,POST",
			"path":         "/say/hello",
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
	data := gateway.Routes{
		Desc:      util.MakeRouteDesc(svariable, endpoint),
		ProtoID:   "00000000000000000359",
		ServiceID: "000022223",
		Endpoint:  endpoint,
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
			return api.NewRequest(opts).Post().SetHeader("Content-Type", "application/json").Resource("route").Body(b)
		},
		Method: "POST",
		URI:    "/route",
		Body:   b.Bytes(),
	}
	headers := map[string]string{"X-API-KEY": conf.MConf().Gateway.Apikey}
	if err := test.Request(reqBody, headers, t); err != nil {
		t.Fatal(err)
	}

}

func aTestWriteRoute(t *testing.T) {
	endpoints := []*registry.Endpoint{
		{
			Name: "KaixinService.Kaixin",
			Metadata: map[string]string{
				"grpc_service": "KaixinService",
				"grpc_method":  "Kaixin",
				"handler":      "HTTP2",
				"method":       "GET,POST",
				"path":         "/say/hello",
			},
		},
		{
			Name: "SayService.Hello",
			Metadata: map[string]string{
				"grpc_service": "SayService",
				"grpc_method":  "Hello",
				"handler":      "HTTP2",
				"method":       "GET,POST",
				"path":         "/say/hello",
			},
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

	data := &task.TaskMsg{
		Action: task.ACTION_CREATE,
		Service: &registry.Service{
			Name:      svariable.Sname,
			Endpoints: endpoints,
		},
		Svariable: svariable,
	}
	plist := []*util.ProtoType{
		{
			PackageName: "aimgo.passport.http2.v1",
			ServiceList: []string{"KaixinService", "WeiwuService"},
			ProtoID:     "00000000000000000359",
		},
		{
			PackageName: "aimgo.passport.http2.v1",
			ServiceList: []string{"SayService"},
			ProtoID:     "00000000000000000390",
		},
	}

	upstreamID := "00000000000000000488"
	r := NewRoute(test.HttpOptions())
	//	根据详情获取service_id
	routeID, err := r.FindByDesc(data, endpoints[0])
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("findServiceIDByDesc[route_id]", routeID)
	}
	fmt.Println(plist)
	//	写入操作
	rID, serr := r.Write( data, true, confPath, upstreamID, plist)
	if serr != nil {
		t.Fatal(serr)
	} else {
		t.Log("update service success [route_id]", rID)
	}
	//	根据service_id获取配置信息
	for _, i := range rID {
		onersp, ferr := r.findByID(i)
		if ferr != nil {
			t.Fatal(ferr)
		} else {
			t.Log(onersp.Node.Key, onersp.Node.Value.Desc)
		}
	}
	//	写入操作
	rid, derr := r.DeleteByUpstreamID(upstreamID)
	if derr != nil {
		t.Fatal(derr)
	} else {
		t.Log("delete service success [route_id]", rid)
	}

}

func TestHttpRoutes(t *testing.T) {

	endpoint := &registry.Endpoint{
		Name: "SayService.Hello",
		Metadata: map[string]string{
			"grpc_service": "SayService",
			"grpc_method":  "Hello",
			"handler":      "HTTP2",
			"method":       "GET,POST",
			"path":         "/say/hello",
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
	data := gateway.Routes{
		Desc:      util.MakeRouteDesc(svariable, endpoint),
		ProtoID:   "00000000000000000359",
		ServiceID: "000022223",
		Endpoint:  endpoint,
		Svariable: svariable,
		Methods:   []string{"GET", "POST"},
	}
	b := new(bytes.Buffer)
	// 数据驱动模板
	err := util.JsonRequestBody(b, "/storage/code/aimgo.config/local/apisix/yaml/web.routes.yaml", "yaml", data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)
}
