package api

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type testcase struct {
	Token  string
	ReqFn  func(opts *Options) *Request
	Method string
	URI    string
	Body   interface{}
	Header map[string]string
	Assert func(req *http.Request) bool
}

var ctx, _ = context.WithTimeout(context.Background(), 1200*time.Millisecond)

var tests = []testcase{
	testcase{
		ReqFn: func(opts *Options) *Request {
			return NewRequest(opts).Get().Resource("service").Params("foo", "bar").Retries(1)
		},
		Method: "GET",
		URI:    "/service?foo=bar",
	},
	testcase{
		ReqFn: func(opts *Options) *Request {
			return NewRequest(opts).Get().Resource("route").Params("foo", "bar").Retries(1)
		},
		Method: "GET",
		URI:    "/route?foo=bar",
	},
	testcase{
		ReqFn: func(opts *Options) *Request {
			return NewRequest(opts).Patch().Resource("service").SetHeader("foo", "bar").Params("foo", "bar")
		},
		Method: "PATCH",
		URI:    "/service?foo=bar",
		Header: map[string]string{"foo": "bar"},
	},
	testcase{
		ReqFn: func(opts *Options) *Request {
			return NewRequest(opts).Patch().Resource("route").SetHeader("foo", "bar").Params("foo", "bar")
		},
		Method: "PATCH",
		URI:    "/route?foo=bar",
		Header: map[string]string{"foo": "bar"},
	},
	testcase{
		ReqFn: func(opts *Options) *Request {
			return NewRequest(opts).Post().Resource("service").Body(map[string]string{"post": "bar"})
		},
		Method: "POST",
		URI:    "/service",
		Body:   map[string]string{"post": "bar"},
	},
	testcase{
		ReqFn: func(opts *Options) *Request {
			return NewRequest(opts).Post().Resource("route").Body(map[string]string{"post": "bar"})
		},
		Method: "POST",
		URI:    "/route",
		Body:   map[string]string{"post": "bar"},
	},
	testcase{
		ReqFn: func(opts *Options) *Request {
			return NewRequest(opts).Put().Resource("service").SubResource("1nssha2idkpq93747nvsw").Body(map[string]string{"bam": "bar"})
		},
		Method: "PUT",
		URI:    "/service/1nssha2idkpq93747nvsw",
		Body:   map[string]string{"bam": "bar"},
	},
	testcase{
		ReqFn: func(opts *Options) *Request {
			return NewRequest(opts).Patch().Resource("service").SubResource("711e52sx291asa1d467131f").Body(map[string]string{"bam": "bar"})
		},
		Method: "PATCH",
		URI:    "/service/711e52sx291asa1d467131f",
		Body:   map[string]string{"bam": "bar"},
	},
	testcase{
		ReqFn: func(opts *Options) *Request {
			return NewRequest(opts).Delete().Resource("service").SubResource("711e52sx291asa1d467131f").Body(map[string]string{"bam": "bar"})
		},
		Method: "DELETE",
		URI:    "/service/711e52sx291asa1d467131f",
		Body:   map[string]string{"bam": "bar"},
	},
	testcase{
		ReqFn: func(opts *Options) *Request {
			return NewRequest(opts).Delete().Context(ctx).Resource("service").SubResource("911e52sx291asa1d467131f").Body(map[string]string{"bam": "bar"})
		},
		Method: "DELETE",
		URI:    "/service/911e52sx291asa1d467131f",
		Body:   map[string]string{"bam": "bar"},
		Token:  "Bearer 32241034405",
	},
}

var wrappedHandler = func(test *testcase, t *testing.T) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(test.Token) > 0 && (len(auth) == 0 || auth != "Bearer "+test.Token) {
			t.Errorf("test case token (%s) did not match expected token (%s)", "Bearer "+test.Token, auth)
		}

		if len(test.Method) > 0 && test.Method != r.Method {
			t.Errorf("test case Method (%s) did not match expected Method (%s)", test.Method, r.Method)
		}

		if len(test.URI) > 0 && test.URI != r.URL.RequestURI() {
			t.Errorf("test case URI (%s) did not match expected URI (%s)", test.URI, r.URL.RequestURI())
		}

		if test.Body != nil {
			var res map[string]string
			decoder := jsoniter.NewDecoder(r.Body)
			if err := decoder.Decode(&res); err != nil {
				t.Errorf("decoding body failed: %v", err)
			}
			if !reflect.DeepEqual(res, test.Body) {
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

func TestRequest(t *testing.T) {
	for _, test := range tests {
		ts := httptest.NewServer(wrappedHandler(&test, t))
		req := test.ReqFn(&Options{
			Host:        ts.URL,
			Client:      &http.Client{},
			BearerToken: &test.Token,
		})
		res := req.Do()
		if res.Error() != nil {
			t.Errorf("request failed with %v", res.Error())
		}
		ts.Close()
	}
}
