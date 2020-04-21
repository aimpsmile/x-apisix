package zaplog

//	执行测试命令

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/micro-in-cn/x-apisix/core/config"
	"github.com/micro-in-cn/x-apisix/core/testconf"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

func initConf() {
	var envType = config.EnvLocal
	if err := testconf.MockConf(testconf.AimEnv(envType)); err != nil {
		log.Panic(err)
	}
	//	覆盖 micro配置日志
	if err := ML().New("test.log.srv.v1", envType, "v1"); err != nil {
		log.Panic(err)
	}
}
func TestInit(t *testing.T) {
	Convey("TestInit", t, func() {
		_, err := testconf.InitEnv(testconf.AimEnv(config.EnvLocal))
		So(err, ShouldBeNil)
		ctx := context.Background()
		traceId := uuid.New().String()
		tctx := ML().NewContext(ctx, String("T-ID", traceId))
		So(tctx.Value(ZapLoggerKey{}), ShouldHaveSameTypeAs, (*zap.Logger)(nil))

		ML().WithContext(tctx).Debug("hahah", Int("num", 23), NamedError("err", fmt.Errorf("this is debug")))
		ML().WithContext(tctx).Info("hahah", Int("num", 24), NamedError("err", fmt.Errorf("this is info")))
		ML().WithContext(tctx).Warn("hahah", Int("num", 25), NamedError("err", fmt.Errorf("this is warning")))
		ML().WithContext(tctx).Error("hahah", Int("num", 26), NamedError("err", fmt.Errorf("this is error")))
	})
}
func TestZaplog(t *testing.T) {
	Convey("TestZaplog", t, func() {
		initConf()
		ML().Debug("hahah", Int("num", 23), NamedError("err", fmt.Errorf("this is debug")))
		ML().Info("hahah", Int("num", 24), NamedError("err", fmt.Errorf("this is info")))
		ML().Warn("hahah", Int("num", 25), NamedError("err", fmt.Errorf("this is warning")))
		ML().Error("hahah", Int("num", 26), NamedError("err", fmt.Errorf("this is error")))
	})
}

func TestTraceContext(t *testing.T) {
	Convey("TestTraceContext", t, func() {
		initConf()
		ctx := context.Background()
		traceId := uuid.New().String()
		tctx := ML().NewContext(ctx, String("T-ID", traceId))
		So(tctx.Value(ZapLoggerKey{}), ShouldHaveSameTypeAs, (*zap.Logger)(nil))

		ML().WithContext(tctx).Debug("hahah", Int("num", 23), NamedError("err", fmt.Errorf("this is debug")))
		ML().WithContext(tctx).Info("hahah", Int("num", 24), NamedError("err", fmt.Errorf("this is info")))
		ML().WithContext(tctx).Warn("hahah", Int("num", 25), NamedError("err", fmt.Errorf("this is warning")))
		ML().WithContext(tctx).Error("hahah", Int("num", 26), NamedError("err", fmt.Errorf("this is error")))
	})
}
