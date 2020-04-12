package middle

import "github.com/micro-in-cn/x-apisix/core/config"

//	检查endpoint是否使用中间件,true:使用、false:不使用
func WhetherUseMiddle(endpoint string) bool {

	if endpoint == config.DebugHealth || endpoint == config.DebugStats {
		return false
	} else {
		return true
	}
}
