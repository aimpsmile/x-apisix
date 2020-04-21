package test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/micro-in-cn/x-apisix/core/lib/api"
	"github.com/micro-in-cn/x-apisix/monitor/conf"
)

type Testcase struct {
	Token  string
	ReqFn  func(opts *api.Options) *api.Request
	Method string
	URI    string
	Body   interface{}
	Header map[string]string
	Assert func(req *http.Request) bool
}

var wrappedHandler = func(test *Testcase, t *testing.T) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(test.Token) > 0 && (auth == "" || auth != "Bearer "+test.Token) {
			t.Errorf("test case token (%s) did not match expected token (%s)", "Bearer "+test.Token, auth)
		}

		if len(test.Method) > 0 && test.Method != r.Method {
			t.Errorf("test case Method (%s) did not match expected Method (%s)", test.Method, r.Method)
		}

		if len(test.URI) > 0 && test.URI != r.URL.RequestURI() {
			t.Errorf("test case URI (%s) did not match expected URI (%s)", test.URI, r.URL.RequestURI())
		}

		if test.Body != nil {
			var res interface{}
			var tres interface{}
			decoder := jsoniter.NewDecoder(r.Body)
			if err := decoder.Decode(&res); err != nil {
				t.Errorf("decoding body failed: %v", err)
			}
			if err := jsoniter.Unmarshal(test.Body.([]byte), &tres); err != nil {
				t.Errorf("decoding body failed: %v", err)
			}
			if !reflect.DeepEqual(res, tres) {
				t.Error("body did not match")
			}
		}

		if test.Header != nil {
			for k, v := range test.Header {
				if r.Header.Get(k) != v {
					t.Error("header did not exist")
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	})
}

//	模拟请求
func Request(test *Testcase, headers map[string]string, t *testing.T) error {
	ts := httptest.NewServer(wrappedHandler(test, t))

	c := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	defer ts.Close()
	req := test.ReqFn(&api.Options{
		Host:        ts.URL,
		Client:      c,
		Headers:     headers,
		BearerToken: &test.Token,
	})
	res := req.Do()
	if res.Error() != nil {
		return res.Error()
	}
	return nil
}

//	模拟配置文件
func MockConf() {

	conf.MockConf(&conf.Conf{
		Filter: []*conf.Filter{
			{
				BU:     "aimgo",
				Stype:  "web",
				Module: "*",
				Ver: map[string]*conf.ModuleVer{
					"*": {
						RouteTpl:    "web.routes.json",
						UpstreamTpl: "web.upstreams.json",
						TplFormat:   "json",
						Hosts:       []string{"web.uqudu.com", "web2.uqudu.com"},
					},
				},
			},
			{
				BU:     "aimgo",
				Stype:  "http2",
				Module: "*",
				Ver: map[string]*conf.ModuleVer{
					"*": {
						RouteTpl:    "grpc.routes.json",
						UpstreamTpl: "grpc.upstreams.json",
						TplFormat:   "json",
						Hosts:       []string{"srv.uqudu.com", "srv.uqudu.com"},
					},
				},
			},
		},
		Gateway: &conf.Gateway{
			Timeout:   5000,
			Retries:   1,
			Apikey:    "edd1c9f034335f136f87ad84b625c8f1",
			Baseurl:   "http://apisix.service:8888/apisix/admin",
			ProtoPath: "/storage/code/aimgo.proto",
			ForbidRoutes: map[string]bool{
				"/stats":  true,
				"/health": true,
			},
		},
		Leader: &conf.Leader{
			Nodes: []string{"etcd.service:2379"},
			Group: "gateway",
		},
		Check: &conf.Check{
			Retries:  3,
			Interval: 60,
		},
	})
}

//  测试环境选项
func HttpOptions() *api.Options {
	MockConf()
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
