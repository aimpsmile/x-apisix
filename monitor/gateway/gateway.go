package gateway

import (
	"github.com/micro-in-cn/x-apisix/core/aimerror"
	"github.com/micro-in-cn/x-apisix/monitor/task"
)

type GatewayI interface {
	Sync(t *task.TaskMsg) (errors aimerror.Errors)
	Cleanup(init bool, Msgs map[string]string) (delSIDs []string, delRIDs []string, errors aimerror.Errors)
	AllDiff(tMsgs []*task.TaskMsg) (diffTasks []*task.TaskMsg, delMsgs map[string]string, err error)
}
type GatewayFunc func() GatewayI
