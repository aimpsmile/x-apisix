package microlog

import (
	"github.com/micro/go-micro/v2/logger"
)

type Options struct {
	logger.Options
}

type namespaceKey struct{}

func WithNamespace(namespace string) logger.Option {
	return logger.SetOption(namespaceKey{}, namespace)
}
