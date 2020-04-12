package uctx

//	针对上下文进行特殊处理的方法库。
import (
	"context"
	"github.com/micro-in-cn/x-apisix/core/lib/ecode"
	"github.com/micro/go-micro/v2/metadata"
	"net/url"
	"strings"
)

//	go-micro上下文件生成
func ConvertContext(ctx context.Context, data map[string][]string) (microCtx context.Context) {
	md := ConvertMetadata(data)

	microCtx = metadata.NewContext(ctx, md)

	return
}

//	map->metadata,
func ConvertMetadata(data map[string][]string) (md metadata.Metadata) {

	md = metadata.Metadata{}
	for k, v := range data {
		md[k] = strings.Join(v, ", ")
	}
	return
}

//  请求信息
func Request(ctx context.Context) (metadata.Metadata, bool) {
	return metadata.FromContext(ctx)
}

//  请求headers的values
func HeaderValues(ctx context.Context, header string) (string, error) {
	h, ok := Request(ctx)
	if !ok {
		return "", ecode.POSTNotFound
	}
	return url.QueryUnescape(h[header])
}

//  http2获取POST参数
func POST(ctx context.Context) (url.Values, error) {
	pstr, err := HeaderValues(ctx, "Micro-Post")
	if err != nil {
		return nil, err
	}
	return url.ParseQuery(pstr)
}

//  http2获取GET参数
func GET(ctx context.Context) (url.Values, error) {

	pstr, err := HeaderValues(ctx, "Micro-Get")
	if err != nil {
		return nil, err
	}
	return url.ParseQuery(pstr)
}

//  http2获取URL参数
func URL(ctx context.Context) (string, error) {
	return HeaderValues(ctx, "Micro-Url")
}

//	获取从网关请求时间
func RequestID(ctx context.Context) (string, error) {
	return HeaderValues(ctx, "X-Request-Id")
}

//	获取ip列表
func ClientIP(ctx context.Context) (ip map[string]string) {
	ip = map[string]string{"local": "", "remote": "", "clientip": "", "proxy": ""}
	if md, ok := Request(ctx); ok {
		ip = map[string]string{"local": md["Local"], "remote": md["Remote"], "clientip": md["X-Real-Ip"], "proxy": md["X-Forwarded-For"]}
	}
	return
}
