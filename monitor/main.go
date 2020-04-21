/*
@Time : 2020/2/12 下午12:38
@Author : songxiuxuan
@File : main.go
@Software: GoLand
*/
package main

import (
	"github.com/micro-in-cn/x-apisix/aimcmd"
	"github.com/micro-in-cn/x-apisix/monitor/cmd"
)

func main() {
	c := cmd.NewCommand()
	c.WithOptions(Options)
	aimcmd.Init(c)
}
