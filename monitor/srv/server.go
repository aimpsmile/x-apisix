// Package 监控服务
package srv

type Server interface {
	Run() error
	CheckAll() error
	String() string
}

// NewMonitor returns a new monitor
func NewServer(opts ...Option) Server {
	return newServer(opts...)
}
