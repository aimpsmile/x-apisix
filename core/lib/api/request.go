package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/micro/go-micro/v2/util/backoff"
)

// Request is used to construct a http request for the k8s API.
type Request struct {
	// the request context
	context context.Context
	client  *http.Client
	header  http.Header
	params  url.Values
	method  string
	host    string
	retries int

	resource    string
	subResource *string
	body        io.Reader

	err error
}

// verb sets method
func (r *Request) verb(method string) *Request {
	r.method = method
	return r
}

func (r *Request) Context(ctx context.Context) *Request {
	r.context = ctx
	return r
}

// 设置重试次数
func (r *Request) Retries(i int) *Request {
	r.retries = i
	return r
}

// Get request
func (r *Request) Get() *Request {
	return r.verb("GET")
}

// Post request
func (r *Request) Post() *Request {
	return r.verb("POST")
}

// Put request
func (r *Request) Put() *Request {
	return r.verb("PUT")
}

// Patch request
func (r *Request) Patch() *Request {
	return r.verb("PATCH")
}

// Delete request
func (r *Request) Delete() *Request {
	return r.verb("DELETE")
}

// Resource is the type of resource the operation is
// for, such as "services", "endpoints" or "pods"
func (r *Request) Resource(s string) *Request {
	r.resource = s
	return r
}

// SubResource sets a subresource on a resource,
// e.g. pods/log for pod logs
func (r *Request) SubResource(s string) *Request {
	r.subResource = &s
	return r
}

// Body pass in a body to set, this is for POST, PUT and PATCH requests
func (r *Request) Body(in interface{}) *Request {
	b := new(bytes.Buffer)
	// if we're not sending YAML request, we encode to JSON
	switch r.header.Get("Content-Type") {
	case "application/yaml", "application/json":
	default:
		if err := jsoniter.NewEncoder(b).Encode(&in); err != nil {
			r.err = err
			return r
		}
		r.body = b
		return r
	}

	// if application/yaml is set, we assume we get a raw bytes so we just copy over
	body, ok := in.(io.Reader)
	if !ok {
		r.err = errors.New("invalid data,not interface io.Reader")
		return r
	}
	// copy over data to the bytes buffer
	if _, err := io.Copy(b, body); err != nil {
		r.err = err
		return r
	}

	r.body = b
	return r
}

// Params isused to set parameters on a request
func (r *Request) Params(key, value string) *Request {
	r.params.Add(key, value)
	return r
}

// SetHeader sets a header on a request with
// a `key` and `value`
func (r *Request) SetHeader(key, value string) *Request {
	r.header.Add(key, value)
	return r
}

// request builds the http.Request from the options
func (r *Request) request() (*http.Request, error) {
	var path string
	path = fmt.Sprintf("%s/%s", r.host, r.resource)
	if r.subResource != nil {
		path += "/" + *r.subResource
	}
	// append any query params
	if len(r.params) > 0 {
		path += "?" + r.params.Encode()
	}

	var req *http.Request
	var err error

	// build request
	if r.context != nil {
		req, err = http.NewRequestWithContext(r.context, r.method, path, r.body)
	} else {
		req, err = http.NewRequest(r.method, path, r.body)
	}
	if err != nil {
		return nil, err
	}

	// set headers on request
	req.Header = r.header
	return req, nil
}

// Do builds and triggers the request
func (r *Request) Do() *Response {
	if r.err != nil {
		return &Response{
			err: r.err,
		}
	}

	req, err := r.request()
	if err != nil {
		return &Response{
			err: err,
		}
	}
	var (
		res  *http.Response
		berr error
	)
	for i := 0; i <= r.retries; i++ {
		t := backoff.Do(i)
		if t.Seconds() > 0 {
			time.Sleep(t)
		}
		// nolint blodyclose
		result, rerr := r.client.Do(req)
		res = result
		berr = rerr
		if rerr == nil {
			break
		}
	}
	if berr != nil {
		return &Response{
			err: berr,
		}
	}

	return newResponse(res, err)
}

// Raw performs a Raw HTTP request to the Kubernetes API
func (r *Request) Raw() (*http.Response, error) {
	req, err := r.request()
	if err != nil {
		return nil, err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type Options struct {
	Host        string
	Resource    string
	Headers     map[string]string
	BearerToken *string
	Client      *http.Client
}

// NewRequest creates a k8s api request
func NewRequest(opts *Options) *Request {
	req := &Request{
		header:   make(http.Header),
		params:   make(url.Values),
		client:   opts.Client,
		resource: opts.Resource,
		host:     opts.Host,
		retries:  0,
	}

	//  设置全局header。主要是认证使用
	for key, val := range opts.Headers {
		req.SetHeader(key, val)
	}

	if opts.BearerToken != nil {
		req.SetHeader("Authorization", "Bearer "+*opts.BearerToken)
	}

	return req
}
