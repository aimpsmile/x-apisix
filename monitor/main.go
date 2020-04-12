/*
@Time : 2020/2/12 下午12:38
@Author : songxiuxuan
@File : main.go
@Software: GoLand
*/
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/micro-in-cn/x-apisix/core/config"
	"github.com/micro-in-cn/x-apisix/core/logger/zaplog"

	"github.com/micro-in-cn/x-apisix/monitor/gateway"
	"github.com/micro-in-cn/x-apisix/monitor/gateway/apisix"

	"github.com/micro-in-cn/x-apisix/aimcmd"
	"github.com/micro-in-cn/x-apisix/core/lib/file"
	"github.com/micro-in-cn/x-apisix/monitor/log"
	"github.com/micro-in-cn/x-apisix/monitor/srv"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	grpcclient "github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	grpctransport "github.com/micro/go-micro/v2/transport/grpc"
	"go.uber.org/zap"
)

const (
	NAME          = "defaultMonitor"
	Usage         = "monitor microservice on service to gateway"
	VERSION       = "v1.0.0"
	GatewayApisix = "apisix"
)

type command struct{}

//  配置服务选项
func withOptions(ctx *cli.Context) []srv.Option {
	address := ctx.String("registry_address")
	confPath := ctx.String("aim_config")
	env := ctx.String("aim_env")
	gw := ctx.String("aim_gateway")

	if address == "" {
		zaplog.ML().Fatal("etcd address not can empty")
	}
	if confPath == "" {
		zaplog.ML().Fatal("config.path not can empty")
	}
	if env == "" {
		zaplog.ML().Fatal("env not can empty")
	}
	if gw == "" {
		zaplog.ML().Fatal("gateway not can empty")
	}
	confPath = fmt.Sprintf("%s/%s/%s", strings.TrimRight(confPath, "/"), env, gw)
	//	check file is exists
	if ok, err := file.Exists(confPath); !ok {
		if err == nil {
			err = fmt.Errorf("file path [%s] not exists", confPath)
		}
		zaplog.ML().Fatal("config.path not exists ", zap.NamedError("error_info", err))
	}

	reg := etcd.NewRegistry(registry.Addrs(address), registry.Timeout(time.Second*3))
	grpcClient := grpcclient.NewClient(
		client.DialTimeout(time.Second*20),
		client.Transport(grpctransport.NewTransport()),
	)
	var defaultGw gateway.GatewayFunc
	switch gw {
	case GatewayApisix:
		defaultGw = func() gateway.GatewayI {
			return apisix.NewClient(gateway.ConfPath(confPath))
		}
	default:
		zaplog.ML().Fatal("gateway config not exists")
	}
	return []srv.Option{
		srv.ConfPath(confPath),
		srv.Registry(reg),
		srv.Client(grpcClient),
		srv.Gateway(defaultGw),
	}

}

//  flags列表
func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "registry_address",
			Value:   "etcd.service:2379",
			Usage:   "Command-separated list of registry addresses",
			EnvVars: []string{"MICRO_REGISTRY_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "aim_config",
			Usage:   "Command-separated list of config path",
			EnvVars: []string{"AIM_CONFIG"},
		},
		&cli.StringFlag{
			Name:    "aim_env",
			Usage:   "Command-separated list of set env",
			EnvVars: []string{"AIM_ENV"},
		},
		&cli.StringFlag{
			Name:    "aim_gateway",
			Value:   GatewayApisix,
			Usage:   "Command-separated list of set gateway type",
			EnvVars: []string{"AIM_GATEWAY"},
		},
	}
}

func checkAllCommand(name, version string) *cli.Command {
	return &cli.Command{
		Name:  "checkall",
		Usage: "Check the all service info",
		Action: func(ctx *cli.Context) error {
			m := srv.NewServer(withOptions(ctx)...)
			if err := log.NewLogger(name, ctx.String("aim_env"), version); err != nil {
				zaplog.ML().Fatal("init logger error", zaplog.NamedError("error_info", err))
			}
			return m.CheckAll()
		},
		Flags: flags(),
	}
}

func (o *command) Name() string {
	name := os.Getenv(config.SnameEnvKey)
	if name == "" {
		return NAME
	} else {
		return name
	}
}
func (o *command) Version() string {
	version := os.Getenv(config.MsverEnvKey)
	if version == "" {
		return VERSION
	} else {
		return version
	}
}
func (o *command) Usage() string {
	return Usage
}

func (o *command) Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "monitor",
		Usage: "Monitor service commanding service",
		Action: func(ctx *cli.Context) error {
			m := srv.NewServer(withOptions(ctx)...)
			if err := log.NewLogger(o.Name(), ctx.String("aim_env"), o.Version()); err != nil {
				zaplog.ML().Fatal("init logger error", zaplog.NamedError("error_info", err))
			}
			return m.Run()
		},
		Flags: flags(),
	}

	return []*cli.Command{command, checkAllCommand(o.Name(), o.Version())}
}

func (o *command) SubFlags() (flags []string) {
	return os.Args
}
func main() {
	aimcmd.Init(new(command))
}
