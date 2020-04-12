package logger

import (
	"io"

	"go.uber.org/zap/zapcore"
)

type WriterI interface {
	Encode(envType string) (encode zapcore.Encoder, whetherSync bool)
	Writer(envType string, level zapcore.Level, opts *Options) (writer io.Writer)
	String() string
}
