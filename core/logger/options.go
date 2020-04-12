package logger

type WriterFunc func() WriterI

// Options 配置项
type Options struct {
	Level string `json:"level"`
	//	日志目录
	LogFileDir string `json:"logFileDir"`
	//	app名称
	AppName string `json:"appName"`
	// 每个日志文件保存的最大尺寸 单位：M
	MaxSize int `json:"maxSize"`
	//	日志文件最多保存多少个备份
	MaxBackups int `json:"maxBackups"`
	// 文件最多保存多少天
	MaxAge int `json:"maxAge"`
	//	获取调用文件信息，跳过中间层
	Skip int `json:"skip"`

	WriterFunc WriterFunc
}

type OptionFunc func(*Options)

//	设置写入对象
func WithWriter(writerFunc WriterFunc) OptionFunc {
	return func(options *Options) {
		options.WriterFunc = writerFunc
	}
}

//	go-micro使用的实例配置
const MICRO = "micro"

//	业务与框架的日志
const LOGIC = "logic"

//	zap日志配置
var ZapConf = map[string][]string{
	MICRO: {"global", "log"},
	LOGIC: {"conf", "log"},
}
