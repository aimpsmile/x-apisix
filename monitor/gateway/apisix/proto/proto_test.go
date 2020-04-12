package proto

import (
	"github.com/micro-in-cn/x-apisix/monitor/conf"
	"github.com/micro-in-cn/x-apisix/monitor/task"
	"github.com/micro-in-cn/x-apisix/monitor/test"
	"github.com/micro/go-micro/v2/registry"
	"testing"
)

const protoPath = "/storage/code/aimgo.proto"

func TestProto(t *testing.T) {

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
			ServiceTpl: "grpc.services.json",
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
	p := NewProto(test.HttpOptions())
	if l, err := p.list(); err != nil {
		t.Fatal(err)
	} else {
		t.Log("get success -->", l.Action)
	}

	//	写入操作
	plist, serr := p.Write(data, protoPath)
	if serr.IsError() {
		t.Fatal(serr.Error())
	} else {
		t.Logf("update service success %+v", plist[0])
	}
	//根据service_id获取配置信息
	onersp, ferr := p.findByID(plist[0].ProtoID)
	if ferr != nil {
		t.Fatal(ferr)
	} else {
		t.Log(onersp.Node.Key, onersp.Node.Value.Content)
	}
	if protoIDs, cerrors := p.ClearList(); cerrors.IsError() {
		t.Fatal(cerrors.Error())
	} else {
		t.Logf("clear all proto success .%+v", protoIDs)
	}
}
