package main

import (
	"time"

	"github.com/micro-in-cn/x-apisix/core/logger/zaplog"
	"github.com/micro-in-cn/x-apisix/monitor/srv"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

func Options(ctx *cli.Context) []srv.Option {
	address := ctx.String("registry_address")
	if address == "" {
		zaplog.ML().Fatal("etcd address not can empty")
	}
	reg := etcd.NewRegistry(registry.Addrs(address), registry.Timeout(time.Second*3))
	return []srv.Option{
		srv.Registry(reg),
	}
}
