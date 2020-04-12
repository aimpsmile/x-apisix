/*
@Time : 2020/2/5 下午4:58
@Author : songxiuxuan
@File : uctx_test.go
@Software: GoLand
*/
package uctx

import (
	"github.com/micro/go-micro/v2/util/addr"
	"net"
	"testing"
)

//  设置服务的metatable
func TestExtractor(t *testing.T) {

	adr := "0.0.0.0"
	host, port, _ := net.SplitHostPort(adr)
	t.Log(host, port)
	adr = host
	t.Log(adr)
	t.Log(addr.Extract(adr))
}
