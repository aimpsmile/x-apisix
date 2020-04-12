package debug

import (
	"net/http"
	_ "net/http/pprof" // #nosec
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	mlog "github.com/micro/go-micro/v2/logger"
)

//	pprof 文件后缀
const ProfSuffix = ".prof"

type profStruct struct {
	cpufile   *os.File
	memfile   *os.File
	blockfile *os.File
	blockname string
}

var pf profStruct

//	打印栈的信息
func Stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}

//	开启http pprof
func StartHTTPProf(httpPort string) {
	if httpPort == "" {
		return
	}
	go func() {
		err := http.ListenAndServe(httpPort, nil)
		mlog.Infof("pprof.http.tracer %v", err)
	}()
}

func TimeCost(name string, start time.Time) {
	terminal := time.Since(start)
	mlog.Debugf("benchmak.time [name]%s[time]%s", name, terminal)
}

//	需要defer 执行方法
func StopProf() {
	stopMem()
	stopCPU()
	stopBlock()
	mlog.Debug("[END_PPROF] success!!!")
}

//	打印cpu与内存的信息
func StartProf(dir, sName, httpPort string) {

	cpuprofile := dir + sName + ".cpu" + ProfSuffix
	memprofile := dir + sName + ".mem" + ProfSuffix
	blockname := sName
	blockfile := dir + sName + ".block" + ProfSuffix

	mlog.Info("pprof.start.tracer")
	if cpuprofile != "" {
		startCPU(cpuprofile)
	}
	if memprofile != "" {
		startMem(memprofile)
	}
	if blockname != "" && blockfile != "" {
		startBlock(blockname, blockfile)
	}
	if httpPort != "" {
		go func() {
			err := http.ListenAndServe(httpPort, nil)
			mlog.Info("pprof.http.tracer %v", err)
		}()
	}
}
func startMem(memprofile string) {
	ff, err := os.Create(memprofile)
	if err != nil {
		mlog.Info("could not create memory profile:  %v", err)
		return
	}
	pf.memfile = ff
}
func stopMem() {
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(pf.memfile); err != nil {
		mlog.Error("could not write memory profile: %v", err)
		return
	}
	pf.memfile.Close()
}
func startCPU(cpuprofile string) {
	f, err := os.Create(cpuprofile)
	if err != nil {
		mlog.Error("could not create CPU profile: %v", err)
		return
	}
	pf.cpufile = f
	if err := pprof.StartCPUProfile(f); err != nil {
		mlog.Error("could not start CPU profile:  %v", err)
		return
	}

}
func stopCPU() {
	pprof.StopCPUProfile()
	pf.cpufile.Close()
}
func startBlock(blockname, blockfile string) {

	p := pprof.NewProfile(blockname)
	if p == nil {
		mlog.Error("could not create block profile %s", blockname)
		return
	}
	f, err := os.Create(blockfile)
	if err != nil {
		mlog.Error("could not create block profile: %v", err)
		return
	}
	pf.blockname = blockname
	pf.blockfile = f
}
func stopBlock() {
	p := pprof.Lookup(pf.blockname)
	if err := p.WriteTo(pf.blockfile, 0); err != nil {
		mlog.Error("could not create memory profile: %v", err)
		return
	}
	pf.blockfile.Close()
}
