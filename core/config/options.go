package config

var (
	GitCommit string
	GitBranch string
	GitTag    string
	BuildDate string
)

//	开发环境常量-列表
const (
	EnvLocal  = "local"
	EnvTest   = "test"
	EnvOnline = "online"
)

//	服务类型
const (
	ServiceAPI   = "api"
	ServiceHTTP2 = "http2"
	ServiceSRV   = "srv"
	ServiceJOB   = "job"
	ServiceWEB   = "web"
	ServiceAPIGW = "apigw"
	ServiceWEBGW = "webgw"
	ServiceTool  = "tool"
)

const (
	//	grpc health endpoint
	DebugHealth = "Debug.Health"
	//	grpc stats endpoint
	DebugStats = "Debug.Stats"
)

//	服务类型列表
var ServiceTypeList = []string{ServiceAPI, ServiceHTTP2, ServiceSRV, ServiceJOB, ServiceWEB, ServiceAPIGW, ServiceWEBGW, ServiceTool}

//	配置中心类型
const (
	RegisterEtcd = "ETCD"
)

//	需要解析的环境变量名称列表
const (
	SnameEnvKey = "MICRO_SERVER_NAME"
	MsverEnvKey = "MICRO_SERVER_VERSION"
	//  服务大版本
	SverEnvKey = "AIM_SVER"
	//  服务类型
	StypeEnvKey = "AIM_STYPE"
	//  服务所属BU
	SbuEnvKey = "AIM_BU"
	//  服务所属模块
	SModuleEnvKey = "AIM_MODULE"
	//  服务环境
	KeyAimEnv = "AIM_ENV"
	//  服务配置路径
	KeyAimConfig = "AIM_CONFIG"
)

//	追踪全局变量
type tracer struct {
	Host string `json:"host"`
}

//	注册中心
type register struct {
	//	注册中心IP:PORT
	Host string `json:"host"`
	//	注册中心类型
	Type string `json:"type"`
	//	超时时间，单位s
	Timeout int `json:"timeout"`
}

//	父级的配置文件
type global struct {
	//	追踪全局变量
	Tracer tracer `json:"tracer"`
	//	注册中心
	Register register `json:"register"`
	//	环境变量
	EnvType string `json:"envType"`
}

//	依赖服务配置列表,
//	解释不配置常量原因：在启动服务的时候，需要做健康检查
type Depend struct {
	Server map[string]string   `json:"server"`
	DB     map[string][]string `json:"db"`
	Cache  map[string][]string `json:"cache"`
	Search map[string][]string `json:"search"`
}

//	全局-服务配置文件
type conf struct {
	//	服务名称
	Sname string `json:"sname"`
	//  bu事业线
	BU string `json:"bu"`
	//  服务版本
	Sver string `json:"sver"`
	//	服务类型
	SType string `json:"stype"`
	//	服务模块
	Module string `json:"module"`
	//	设置pprof环境变量:[http：port：6060,非空:pprof]
	Sprof string `json:"sprof"`
	//	是否关闭信号
	Signal bool `json:"signal"`
	//  服务的metatable
	Smetatable map[string]string `json:"smetatable"`
	//	依赖服务
	Depend Depend `json:"depend"`
}

//	全局配置文件
type Gconf struct {
	conf
	global
}

func (c *Gconf) String() string {
	return "core.config"
}

//  设置服务的metatable
func ServerMetatable(meta map[string]string) map[string]string {
	if len(meta) == 0 {
		meta = map[string]string{}
	}
	meta["tag"] = GitTag
	meta["branch"] = GitBranch
	meta["commit"] = GitCommit
	meta["time"] = BuildDate
	return meta
}
