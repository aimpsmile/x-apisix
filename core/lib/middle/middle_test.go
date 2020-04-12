package middle

import (
	"testing"

	"github.com/micro-in-cn/x-apisix/core/config"
)

func TestWhetherUseMiddle(t *testing.T) {
	if WhetherUseMiddle(config.DebugHealth) {
		t.Fatal("health can't use middle")
	}
	if WhetherUseMiddle(config.DebugStats) {
		t.Fatal("stats can't use middle")
	}
}
