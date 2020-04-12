package log

import (
	"github.com/micro-in-cn/x-apisix/core/logger/microlog"
	"github.com/micro-in-cn/x-apisix/core/logger/zaplog"
	"github.com/micro/go-micro/v2/logger"
)

func NewLogger(sName, aimEnv, sver string) error {
	//	初始化micro日志
	if err := zaplog.ML().New(sName, aimEnv, sver); err != nil {
		return err
	}
	l, err := microlog.NewLogger(microlog.WithNamespace("go-micro"))
	if err != nil {
		return err
	}
	logger.DefaultLogger = l
	return nil
}
func ML() *zaplog.ZapLogger {
	return zaplog.ML()
}
