package writer

import (
	"fmt"
	"github.com/micro-in-cn/x-apisix/core/config"
	"github.com/micro-in-cn/x-apisix/core/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
)

type Writer struct {
}

//	日志消息格式
func (w *Writer) Encode(envType string) (encode zapcore.Encoder, whetherSync bool) {
	switch envType {
	case config.EnvLocal:
		encode = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		whetherSync = false
	default:
		encode = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		whetherSync = true
	}
	return
}

//	写入配置项
func (w *Writer) lumberjack(level zapcore.Level, opts *logger.Options) io.Writer {
	sp := string(filepath.Separator)
	fN := fmt.Sprintf("%s.log", level.String())
	return &lumberjack.Logger{
		Filename:   opts.LogFileDir + sp + opts.AppName + "-" + fN,
		MaxSize:    opts.MaxSize,
		MaxBackups: opts.MaxBackups,
		MaxAge:     opts.MaxAge,
		Compress:   true,
		LocalTime:  true,
	}
}

//	外部写入日志
func (w *Writer) Writer(envType string, level zapcore.Level, opts *logger.Options) (writer io.Writer) {
	switch envType {
	case config.EnvLocal:
		writer = os.Stdout
	default:
		writer = w.lumberjack(level, opts)
	}
	return
}
func (w *Writer) String() string {
	return "core.log.zaplog.writer"
}

func GetWriter() logger.WriterI {
	return &Writer{}
}
