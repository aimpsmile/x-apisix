package gateway

import "strings"

type Options struct {
	ConfPath string
}
type Option func(*Options)

//  网关-配置文件路径
func ConfPath(confPath string) Option {
	return func(options *Options) {
		options.ConfPath = strings.TrimRight(confPath, "/")
	}
}
