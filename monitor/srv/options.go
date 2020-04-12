package srv

import (
	"github.com/micro-in-cn/x-apisix/monitor/gateway"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
)

type Options struct {
	Client   client.Client
	Registry registry.Registry
	Gateway  func() gateway.GatewayI
	ConfPath string
}

type Option func(*Options)

func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func Gateway(rfunc gateway.GatewayFunc) Option {
	return func(o *Options) {
		o.Gateway = rfunc
	}
}

func ConfPath(confPath string) Option {
	return func(options *Options) {
		options.ConfPath = confPath
	}

}
