package zaplog

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"github.com/micro-in-cn/x-apisix/core/logger"
	"github.com/micro-in-cn/x-apisix/core/logger/zaplog/writer"

	"github.com/micro-in-cn/x-apisix/core/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const sp = string(filepath.Separator)

type ZapLoggerKey struct{}

//	方便其它应用指定日志类型
const (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel
)

type ZapLogger struct {
	*zap.Logger
	//	配置名称
	confName string
	//	是否同步
	whetherSync bool
	//	日志写入对象
	w logger.WriterI
	//	log.config
	opts *logger.Options
	//	级别
	level zapcore.Level
	//	锁
	mtx *sync.RWMutex
}

var (
	logiclog = defaultLog(logger.LOGIC, "aim.logic")
	microLog = defaultLog(logger.MICRO, "aim.micro")
)

func defaultLog(confName, sname string) (zlog *ZapLogger) {

	envType := config.EnvLocal
	zapConfig := zap.NewDevelopmentConfig()
	LogFileDir, _ := filepath.Abs(filepath.Dir(filepath.Join(".")))
	LogFileDir += sp + "logs" + sp
	logOpts := &logger.Options{
		//	日志名称
		AppName: sname,
		// 默认输出到程序运行目录的logs子目录
		LogFileDir: LogFileDir,
		//	日志文件多大，进行切割
		MaxSize: 50,
		//	备份的份数
		MaxBackups: 3,
		//	旧文件保留的day
		MaxAge:     30,
		WriterFunc: writer.GetWriter,
	}
	zlog = &ZapLogger{
		Logger:      nil,
		confName:    confName,
		opts:        logOpts,
		w:           logOpts.WriterFunc(),
		whetherSync: false,
		level:       zapConfig.Level.Level(),
		mtx:         new(sync.RWMutex),
	}
	opts, whetherSync := zlog.zapWriter(envType, "")
	zlog.whetherSync = whetherSync
	zlog.Logger, _ = zapConfig.Build(opts, zap.AddStacktrace(ErrorLevel))

	return zlog
}

/*
让go-micro的日志，转入到本地日志系统
systemLog := ML().New("testlog", envType)
mlog.SetZapLogger(systemLog)
systemLog.Info("this is kaixin logic")
*/
func (l *ZapLogger) New(sName, envType, sver string, opts ...logger.OptionFunc) error {
	for _, o := range opts {
		o(l.opts)
	}
	err := l.initZapCfg(sName, envType, sver)

	if err != nil {
		return err
	}
	l.Info("[initZapLogger] zap plugin init success.", String("completed", l.confName))
	return nil
}

//	框架日志使用
func LL() *ZapLogger {
	return logiclog
}

//	go-micro日志使用
func ML() *ZapLogger {
	return microLog
}

//	加载业务的配置文件
func (l *ZapLogger) loadCfg(sName string) (err error) {
	if _, ok := logger.ZapConf[l.confName]; !ok {
		err = fmt.Errorf("配置项%v不存在", l.confName)
		return
	}
	err = config.C().StructList(&l.opts, logger.ZapConf[l.confName]...)
	if err != nil {
		return
	}
	if l.opts == nil {
		err = fmt.Errorf("加载配置项%v失败", l.confName)
		return
	}
	l.opts.AppName = sName
	return
}

//	zap.writer根据level确定写入路径
func (l *ZapLogger) zapWriter(envType, sver string) (zap.Option, bool) {

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel && lvl-l.level > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && lvl-l.level > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && lvl-l.level > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && lvl-l.level > -1
	})

	encode, whetherSync := l.w.Encode(envType)
	//	设置服务名称
	encode.AddString("sname", l.opts.AppName)
	//  设置服务版本
	encode.AddString("sver", sver)
	// 日志文件
	cores := []zapcore.Core{
		//	error 及以上
		zapcore.NewCore(encode, zapcore.AddSync(l.w.Writer(envType, zapcore.ErrorLevel, l.opts)), errPriority),

		//	warn
		zapcore.NewCore(encode, zapcore.AddSync(l.w.Writer(envType, zapcore.WarnLevel, l.opts)), warnPriority),

		//	info
		zapcore.NewCore(encode, zapcore.AddSync(l.w.Writer(envType, zapcore.InfoLevel, l.opts)), infoPriority),

		//	debug
		zapcore.NewCore(encode, zapcore.AddSync(l.w.Writer(envType, zapcore.DebugLevel, l.opts)), debugPriority),
	}

	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	}), whetherSync
}

//	初始化zap日志的配置文件
func (l *ZapLogger) initZapCfg(sName, envType, sver string) (err error) {

	//	根据配置加载日志配置文件
	err = l.loadCfg(sName)
	if err != nil {
		return
	}
	//	设置服务级别
	err = l.level.Set(l.opts.Level)
	if err != nil {
		return
	}
	opts, whetherSync := l.zapWriter(envType, sver)
	l.whetherSync = whetherSync
	l.mtx.Lock()
	l.Logger = l.WithOptions(opts, zap.AddCallerSkip(l.opts.Skip), zap.AddStacktrace(ErrorLevel))
	l.mtx.Unlock()

	return
}

// 为日志上下文，增加字段
func (l *ZapLogger) NewContext(ctx context.Context, fields ...zapcore.Field) context.Context {
	return context.WithValue(ctx, ZapLoggerKey{}, l.WithContext(ctx).With(fields...))
}

// 设置上下文选项获取zap对象
func (l *ZapLogger) WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return l.Logger
	}

	if ctxZapLogger, ok := ctx.Value(ZapLoggerKey{}).(*zap.Logger); ok {
		return ctxZapLogger
	}
	return l.Logger
}

//	当前系统配置日志级别
func (l *ZapLogger) Level() zapcore.Level {
	return l.level
}

//	设置日志格式
func (l *ZapLogger) Encode(envType string) (encode zapcore.Encoder) {
	encode, _ = l.w.Encode(envType)
	return encode
}

//	设置写入对象
func (l *ZapLogger) Writer(envType string) (w io.Writer) {
	return l.w.Writer(envType, l.Level(), l.opts)
}

//	针对zap.Sync() 进行包装
func (l *ZapLogger) Report() (err error) {
	if l.whetherSync {
		err = l.Logger.Sync()
	}
	return
}

func (l *ZapLogger) String() string {
	return "core.zaplog"
}
