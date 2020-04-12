package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// Errors ...
var (
	ErrNotFound = errors.New("http.client: resource not found")
	ErrDecode   = errors.New("http.client: error decoding")
	ErrUnknown  = errors.New("http.client: unknown error")
)

// Status is an object that is returned when a request
// failed or delete succeeded.
type Status struct {
	Kind    string `json:"kind"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Code    int    `json:"code"`
}

type Response struct {
	res *http.Response
	err error
}

// Error returns an error
func (r *Response) Error() error {
	return r.err
}

// StatusCode returns status code for response
func (r *Response) StatusCode() int {
	return r.res.StatusCode
}

// Into decode body into `data`
func (r *Response) Into(data interface{}) error {
	if r.err != nil {
		return r.err
	}

	defer func() {
		if err := r.res.Body.Close(); err != nil {
			fmt.Printf("http close body is error %v .", err)
		}
	}()
	decoder := jsoniter.NewDecoder(r.res.Body)
	err := decoder.Decode(&data)
	if err != nil {
		return ErrDecode
	}

	return r.err
}

func (r *Response) Close() error {
	return r.res.Body.Close()
}

func newResponse(res *http.Response, err error) *Response {
	r := &Response{
		res: res,
		err: err,
	}

	if err != nil {
		return r
	}

	if r.res.StatusCode == http.StatusOK ||
		r.res.StatusCode == http.StatusCreated ||
		r.res.StatusCode == http.StatusNoContent {
		// Non error status code
		return r
	}

	if r.res.StatusCode == http.StatusNotFound {
		r.err = ErrNotFound
		return r
	}

	b, err := ioutil.ReadAll(r.res.Body)
	if err == nil {
		r.err = errors.New(string(b))
		return r
	}

	r.err = ErrUnknown

	return r
}
