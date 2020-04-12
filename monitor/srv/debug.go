package srv

import (
	"context"
	"net/http"
	"time"

	"fmt"

	"github.com/micro-in-cn/x-apisix/core/lib/api"
	"github.com/micro/go-micro/v2/client"
	grpcclient "github.com/micro/go-micro/v2/client/grpc"
	pb "github.com/micro/go-micro/v2/debug/service/proto"
	"github.com/micro/go-micro/v2/registry"
	grpctransport "github.com/micro/go-micro/v2/transport/grpc"
)

const (
	TransportGrpc = "grpc"
	TransportHttp = "http"
)

func grcpClient() client.Client {
	return grpcclient.NewClient(
		client.DialTimeout(time.Second*20),
		client.Transport(grpctransport.NewTransport()),
	)
}
func httpClient(host string) *api.Options {
	c := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	return &api.Options{
		Host:   fmt.Sprintf("http://%s", host),
		Client: c,
	}
}

func Health(name string, node *registry.Node, retrys int) (*pb.HealthResponse, error) {
	// check the transport matches
	switch node.Metadata["transport"] {
	case TransportGrpc:
		debug := pb.NewDebugService(name, grcpClient())
		return debug.Health(
			context.Background(),
			// empty health request
			&pb.HealthRequest{},
			// call this specific node
			client.WithAddress(node.Address),
			// retry in the event of failure
			client.WithRetries(retrys),
		)
	default:
		rsp := &pb.HealthResponse{}
		err := api.NewRequest(httpClient(node.Address)).
			Get().
			SetHeader("Content-Type", "application/json").
			Resource("health").
			Retries(retrys).
			Do().
			Into(rsp)
		return rsp, err
	}
}
func Stats(name string, node *registry.Node, retrys int) (*pb.StatsResponse, error) {
	switch node.Metadata["transport"] {
	case TransportGrpc:
		debug := pb.NewDebugService(name, grcpClient())
		return debug.Stats(
			context.Background(),
			// empty health request
			&pb.StatsRequest{},
			// call this specific node
			client.WithAddress(node.Address),
			// retry in the event of failure
			client.WithRetries(retrys),
		)

	default:
		rsp := &pb.StatsResponse{}
		err := api.NewRequest(httpClient(node.Address)).
			Get().
			SetHeader("Content-Type", "application/json").
			Retries(retrys).
			Resource("stats").
			Do().
			Into(rsp)
		return rsp, err
	}
}
