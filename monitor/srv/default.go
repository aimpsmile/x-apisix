package srv

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/micro/go-micro/v2/sync/leader"
	"github.com/micro/go-micro/v2/sync/leader/etcd"

	"github.com/micro-in-cn/x-apisix/core/logger/zaplog"
	"github.com/micro-in-cn/x-apisix/monitor/conf"
	"github.com/micro-in-cn/x-apisix/monitor/gateway"
	"github.com/micro-in-cn/x-apisix/monitor/gateway/apisix"
	"github.com/micro-in-cn/x-apisix/monitor/task"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/util/backoff"
)

type server struct {
	options Options

	registry registry.Registry
	client   client.Client
	gateway  gateway.GatewayI

	sync.RWMutex
	isclear      int32
	numThread    int
	running      bool
	exit         chan error
	closeConsume chan bool
	jobChan      chan *task.TaskMsg
	errChan      chan error
	wg           *sync.WaitGroup
	ele          leader.Elected
}

//	监听cmd命令信号
func listenSign() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		sign := <-ch
		zaplog.ML().Info("listen sign to stop server.", zaplog.String("sign", sign.String()))
		cancel()
	}()
	return ctx, cancel
}

//	推送消息到队列中
func (m *server) pushJob(msg *task.TaskMsg) {
	m.jobChan <- msg
}

func (m *server) popJob() (msg *task.TaskMsg) {
	msg = <-m.jobChan
	return
}

//	获取任务队列数量
func (m *server) lenJob() int {
	return len(m.jobChan)
}

//	判断服务处理的队列
func (m *server) closeJob() {
	close(m.jobChan)
}

//	主动健康检查服务运行状态
func (m *server) health(service *registry.Service) {
	j := 0
	for _, node := range service.Nodes {

		rsp, err := Health(service.Name, node, conf.MConf().Check.Retries)
		if err != nil || rsp.Status != "ok" {
			zaplog.ML().Error("[HEALTH.DELETE.NODE]",
				zaplog.String("snmae", service.Name),
				zaplog.String("ID", node.Id),
				zaplog.String("address", node.Address),
				zaplog.NamedError("error_info", err),
			)
			continue
		}
		service.Nodes[j] = node
		j++
	}
	service.Nodes = service.Nodes[:j]
}

//	创建task任务
func (m *server) newTask(action string, service *registry.Service, init bool) *task.TaskMsg {
	t := task.NewMsg(action, service)
	if t == nil {
		return nil
	}
	if t.Action != task.ACTION_DELETE {
		//	入队列之前做健康检查
		m.health(t.Service)
		if len(t.Service.Nodes) == 0 {
			zaplog.ML().Warn("[SERVER.EXCEPTION]set action to delete as service node not exists ",
				zaplog.Int("node_num", len(t.Service.Nodes)),
				zaplog.String("snmae", t.Svariable.Sname),
				zaplog.String("version", t.Svariable.Version),
			)
			//	切换成delete状态
			t.Action = task.ACTION_DELETE
		}
	}

	if init {
		m.pushJob(t)
	}

	return t
}

//	判断服务是否退出
func (m *server) quit(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

//	watch服务变更
func (m *server) watch(ctx context.Context) {
	var a int
	var watcher registry.Watcher

	for {
		if m.quit(ctx) {
			break
		}
		stop := make(chan bool)
		//	注册中心监听
		w, rerr := m.registry.Watch()
		if rerr != nil {
			d := backoff.Do(a)
			if a > 3 {
				zaplog.ML().Error("[WATCH.ERROR]",
					zaplog.NamedError("error_info", rerr),
					zaplog.Int("retrers", a))
				a = 0
			}
			time.Sleep(d)
			a++
			continue
		}
		a = 0
		watcher = w

		//	停止watcher
		go func() {
			defer watcher.Stop()
			select {
			case <-stop:
				return
			case <-ctx.Done():
				return
			}
		}()

		for {

			s, werr := watcher.Next()
			if werr != nil {
				zaplog.ML().Error("[WATCH.NEXT]", zaplog.NamedError("error_info", werr))
				close(stop)
				break
			}

			if s != nil {
				//	检测是否符合网关需要的规则
				_, err := task.ShouldSnameBeConf(s.Service.Name, s.Service.Version, conf.MConf())
				if err != nil {
					continue
				}
				//	重新拉取注册中心的配置
				servcieList, err := m.registry.GetService(s.Service.Name)

				//	确实被删除了
				if err == registry.ErrNotFound {
					//	很有可能不存在或者连接失败，不管了先按照watch进行同步
					m.newTask(s.Action, s.Service, true)
					continue
				}
				//	获取数据发生了错误，需要记录日志
				if err != nil {
					zaplog.ML().Error("get.service.list by sname", zaplog.NamedError("error_info", err))
					continue
				}
				for _, service := range servcieList {
					m.newTask(task.ACTION_UPDATE, service, true)
				}
			}
		}
	}
}

//	检查服务运行的状态
func (m *server) check(service *registry.Service) {
	for _, node := range service.Nodes {

		rsp, err := Health(service.Name, node, conf.MConf().Check.Retries)
		if err != nil {
			zaplog.ML().Error("[CHECK.SERVER.HEALTH]",
				zaplog.String("snmae", service.Name),
				zaplog.Reflect("node", node),
				zaplog.NamedError("error_info", err))
		} else if rsp != nil {
			zaplog.ML().Info("[CHECK.SERVER.HEALTH]",
				zaplog.String("snmae", service.Name),
				zaplog.Reflect("node", node),
				zaplog.String("status", rsp.Status))
		}
		sresp, serr := Stats(service.Name, node, conf.MConf().Check.Retries)
		if serr != nil {
			zaplog.ML().Error("[CHECK.SERVER.STATS]",
				zaplog.String("snmae", service.Name),
				zaplog.Reflect("node", node),
				zaplog.NamedError("error_info", serr))
		} else {
			zaplog.ML().Info("[CHECK.SERVER.STATS]",
				zaplog.String("snmae", service.Name),
				zaplog.Reflect("node", node),
				zaplog.NamedError("error_info", serr))
			zaplog.ML().Info("[CHECK.SERVER.STATS]",
				zaplog.Uint64("time", sresp.Timestamp),
				zaplog.String("snmae", service.Name),
				zaplog.Reflect("node", node),
				zaplog.Uint64("starttime", sresp.Started),
				zaplog.Uint64("uptime", sresp.Uptime),
				zaplog.Uint64("reqs", sresp.Requests),
				zaplog.Uint64("threads", sresp.Threads),
				zaplog.Uint64("gc", sresp.Gc),
				zaplog.Uint64("memory", sresp.Memory),
				zaplog.Uint64("errors", sresp.Errors),
			)
		}
	}
	m.wg.Done()
}

func (m *server) all(init bool) error {
	services, err := m.registry.ListServices()
	if err != nil {
		return err
	}
	//  全量进行同步
	var taskMsgs []*task.TaskMsg
	for _, s := range services {
		t := m.newTask(task.ACTION_UPDATE, s, init)
		if t == nil {
			continue
		}
		taskMsgs = append(taskMsgs, t)
	}

	for {
		if m.lenJob() > 0 {
			continue
		} else {
			atomic.AddInt32(&m.isclear, 1)
			break
		}
	}
	diffTasks, delMsgs, diffErr := m.gateway.AllDiff(taskMsgs)
	zaplog.ML().Info("[diff]gateway diff service",
		zaplog.Reflect("diff_task", diffTasks),
		zaplog.Reflect("del_msg", delMsgs),
		zaplog.NamedError("error_infos", diffErr))
	if diffErr != nil {
		atomic.AddInt32(&m.isclear, -1)
		return err
	}
	for _, dt := range diffTasks {
		m.pushJob(dt)
	}
	delSIDs, delRIDs, errors := m.gateway.Cleanup(init, delMsgs)
	zaplog.ML().Info("[cleanup]gateway cleanup.service.info",
		zaplog.Strings("del_upstream_list", delSIDs),
		zaplog.Strings("del_route_list", delRIDs),
		zaplog.Errors("error_infos", errors))
	atomic.AddInt32(&m.isclear, -1)
	return nil
}

//	生产者-监听注册中心服务列表的变更
func (m *server) product(ctx context.Context) {
	err := m.all(true)
	if err != nil {
		//	关闭任务通道
		m.closeJob()
		m.errChan <- err
		return
	}
	//	监听注册中心服务变化
	m.watch(ctx)
	//	关闭任务通道
	m.closeJob()
}

//	消费者-统一处理服务与其它服务进行通信，总控
func (m *server) consume(ctx context.Context) {
	for t := range m.jobChan {
		for atomic.LoadInt32(&m.isclear) > 0 {
		}
		syncErrs := m.gateway.Sync(t)
		if syncErrs.IsError() {
			zaplog.ML().Error("sync is errors ",
				zaplog.Int("retries", t.Retries),
				zaplog.String("actiotn", t.Action),
				zaplog.String("snmae", t.Service.Name),
				zaplog.Errors("error_infos", syncErrs))
			//	一种错误只能重试三次
			if t.Retries < conf.MConf().Check.Retries {
				zaplog.ML().Info("retry push task to job", zaplog.Int("num", t.Retries))
				t.Retries += 1
				m.jobChan <- t
			}
		}
	}
	m.wg.Wait()

	<-ctx.Done()
	m.closeConsume <- true
}

//	运行具体任务-控制
func (m *server) run(ctx context.Context) {
	t := time.NewTicker(time.Second * time.Duration(conf.MConf().Check.Interval))

	//	增量同步watch
	go m.product(ctx)
	//	消费服务
	go m.consume(ctx)

Loop:
	for {
		select {
		//	服务信号stop
		case <-m.closeConsume:
			m.ele.Revoked()
			break Loop
		case <-t.C:
			if m.lenJob() > 0 {
				continue
			}
			err := m.all(false)
			if err != nil {
				zaplog.ML().Error("an error occurred with the cleanup service",
					zaplog.String("leader_id", m.ele.Id()),
					zaplog.NamedError("error_info", err))
			}
		}
	}
	if err := m.ele.Resign(); err != nil {
		zaplog.ML().Error("Resign leader election error occurred",
			zaplog.String("leader_id", m.ele.Id()),
			zaplog.NamedError("error_info", err))
	}
	//	任务全部结束
	m.exit <- nil
}

//	开启服务
func (m *server) start(ctx context.Context) error {
	lconf := conf.MConf().Leader
	eleader := etcd.NewLeader(leader.Nodes(lconf.Nodes...), leader.Group(lconf.Group))
	if ele, err := eleader.Elect(lconf.ID); err != nil {
		return err
	} else {
		m.ele = ele
	}

	m.Lock()
	defer m.Unlock()

	if m.running {
		return nil
	}
	//	生产服务
	go m.run(ctx)
	m.running = true
	return nil

}

//	停止运行服务
func (m *server) stop() error {
	m.Lock()
	defer m.Unlock()
	if !m.running {
		zaplog.ML().Info("server already stop")
		return nil
	}
	//	服务已经退出
	err := <-m.exit
	//	退出完毕，退出过程中出现错误，进行记录
	m.running = false

	return err
}

//	控制主流程
func (m *server) Run() error {

	ctx, cancel := listenSign()
	// create a new monitor
	if err := m.start(ctx); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		zaplog.ML().Info("because the sign close server ing...")
	case err := <-m.errChan:
		cancel()
		zaplog.ML().Error("because the error close server ing...", zaplog.NamedError("error_info", err))
	}
	// stop monitor
	if err := m.stop(); err != nil {
		return err
	}
	return nil
}

//	检查所有服务列表
func (m *server) CheckAll() error {
	services, err := m.registry.ListServices()
	if err != nil {
		zaplog.ML().Error("[CHECK.ERROR] ", zaplog.NamedError("error_info", err))
		return nil
	}
	for _, service := range services {
		m.wg.Add(1)
		go m.check(service)
	}
	m.wg.Wait()
	return nil
}

func (m *server) String() string {
	return "monitor.srv"
}

func newServer(opts ...Option) Server {
	// 发生宕机时，获取panic传递的上下文并打印
	defer func() {
		if r := recover(); r != nil {
			zaplog.ML().Fatal("srv.panic", zaplog.Any("panic_error_info", r))
		}
	}()

	options := Options{
		Client:   client.DefaultClient,
		Registry: registry.DefaultRegistry,
	}

	for _, o := range opts {
		o(&options)
	}
	//  初始化配置文件
	confFile := fmt.Sprintf("%s/%s", options.ConfPath, "conf.yml")
	if err := conf.InitConf(confFile); err != nil {
		zaplog.ML().Fatal("service init conf error ", zaplog.String("conf_path", confFile), zaplog.NamedError("error_info", err))
	}
	//	gateway配置目录与monitor配置相同文件
	var gw gateway.GatewayI
	if options.Gateway != nil {
		gw = options.Gateway()
	} else {
		gw = apisix.NewClient(gateway.ConfPath(options.ConfPath))
	}

	numThread := runtime.NumCPU() * 2

	return &server{
		options:      options,
		numThread:    numThread,
		isclear:      0,
		wg:           &sync.WaitGroup{},
		exit:         make(chan error),
		closeConsume: make(chan bool),
		jobChan:      make(chan *task.TaskMsg, numThread),
		errChan:      make(chan error),
		client:       options.Client,
		registry:     options.Registry,
		gateway:      gw,
	}
}
