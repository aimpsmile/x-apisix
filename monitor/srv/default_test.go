package srv

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/micro-in-cn/x-apisix/core/logger/zaplog"

	"github.com/micro-in-cn/x-apisix/monitor/gateway"
	"github.com/micro-in-cn/x-apisix/monitor/gateway/apisix"
	"go.uber.org/zap"

	"github.com/micro-in-cn/x-apisix/monitor/log"
	"github.com/micro/go-micro/v2/client"
	grpcclient "github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	grpctransport "github.com/micro/go-micro/v2/transport/grpc"
)

const (
	ENV     = "local"
	NAME    = "aimGoMonitor"
	VERSION = "0.5.1"
)

func withOptions() []Option {
	reg := etcd.NewRegistry(registry.Addrs("etcd.service:2379"), registry.Timeout(time.Second*3))
	grpcClient := grpcclient.NewClient(
		client.DialTimeout(time.Second*20),
		client.Transport(grpctransport.NewTransport()),
	)

	confPath := "/storage/code/aimgo.config"
	env := "local"
	gw := "apisix"
	confPath = fmt.Sprintf("%s/%s/%s", strings.TrimRight(confPath, "/"), env, gw)
	defaultGw := func() gateway.GatewayI {
		return apisix.NewClient(gateway.ConfPath(confPath))
	}
	return []Option{
		ConfPath(confPath),
		Registry(reg),
		Client(grpcClient),
		Gateway(defaultGw),
	}
}

func aTestCheck(t *testing.T) {
	m := NewServer(withOptions()...)
	if err := log.NewLogger(NAME, ENV, VERSION); err != nil {
		zaplog.ML().Fatal("init logger error", zaplog.NamedError("error_info", err))
	}
	err := m.CheckAll()
	if err != nil {
		log.ML().Fatal("启动服务失败", zap.NamedError("error_info", err))
	}
}

func TestMonitor(t *testing.T) {
	m := NewServer(withOptions()...)
	if err := log.NewLogger(NAME, ENV, VERSION); err != nil {
		zaplog.ML().Fatal("init logger error", zaplog.NamedError("error_info", err))
	}
	err := m.Run()
	if err != nil {
		log.ML().Fatal("启动服务失败", zap.NamedError("error_info", err))
	}
}
